package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

/* 投票的集中情况
direction = 1 时, 有两种情况
    1. 之前没有投过票, 现在投赞成票 --> 更新分数和投票记录    差值的绝对值: 1  +432
    2. 之前投反对票, 现在改投赞成票 --> 更新分数和投票记录    差值的绝对值: 2  +432*2
direction = 0 时, 有两种情况
    1. 之前投赞成票, 现在要取消投票 --> 更新分数和投票记录    差值的绝对值: 1  -432
    2. 之前投反对票, 现在要取消投票 --> 更新分数和投票记录    差值的绝对值: 1  +432
direction = 1 时, 有两种情况
    1. 之前没有投过票, 现在投反对票 --> 更新分数和投票记录    差值的绝对值: 1  -432
    2. 之前投赞成票, 现在改投反对票 --> 更新分数和投票记录    差值的绝对值: 2  -432*2

投票的限制:
每个帖子自发表之日起一个星期之内允许用户投票, 超过一个星期就不允许投票了
    1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
    2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {

	// 事务
	pipeline := client.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()

	return err
}

func VoteForPost(userID, postID string, value float64) error {

	// 1. 判断投票规则
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2. 更新帖子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	// 如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //  计算两次投票的差值

	// 2和3需要放到一个pipeline事务中操作
	pipeline := client.TxPipeline()

	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录该用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}

	_, err := pipeline.Exec()
	return err
}
