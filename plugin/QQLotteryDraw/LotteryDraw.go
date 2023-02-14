package QQLotteryDraw

import (
	"fmt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var LDLineUp = make(map[string][]int64)

func init() {
	engine := control.Register("QQ群抽奖", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "发送“抽奖帮助”获取抽奖使用帮助",
		OnEnable: func(ctx *zero.Ctx) {
			process.SleepAbout1sTo2s()
			ctx.SendChain(message.Text("开抽！~（抓起小皮鞭）"))
		},
		OnDisable: func(ctx *zero.Ctx) {
			process.SleepAbout1sTo2s()
			ctx.SendChain(message.Text("没得抽了~(收起小皮鞭)"))
		},
	})
	var LDConfig = make([]string, 3)
	engine.OnPrefixGroup([]string{"开启抽奖", "开始抽奖"}, zero.AdminPermission, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		LDtxt := ctx.State["args"].(string)
		LDConfig = strings.Split(LDtxt, " ")
		NumberOPeople, NOPerr := strconv.Atoi(LDConfig[0])
		TimeCache := strings.TrimSpace(LDConfig[1])
		LDTime, TimeErr := strconv.Atoi(TimeCache)
		fmt.Println("时间", LDTime)
		if len(LDConfig) != 3 && NOPerr != nil && TimeErr != nil && LDTime > 0 {
			ctx.SendChain(message.Text("参数数量不正确"))
		} else {
			_, IsEx := LDLineUp[LDConfig[2]]
			if IsEx {
				ctx.SendChain(message.Text("当前抽奖关键词已存在，请更换关键词创建抽奖"))
			} else {
				LDLineUp[LDConfig[2]] = []int64{}
				ctx.SendChain(message.Text("抽奖开始了\n本次在" + LDConfig[1] + "秒内共抽取" + LDConfig[0] + "个人\n赶快发送\n\n参加抽奖 " + LDConfig[2] + "\n\n来参与抽奖吧"))
				fmt.Println(time.Now().Unix())
				go func() {
					LDKeyString := LDConfig[2]
					<-time.After(time.Duration(LDTime) * time.Second)
					fmt.Println(time.Now().Unix())
					fmt.Println(NumberOPeople)
					AllUserId := LDLineUp[LDKeyString]
					var WinUser []int64
					LdNumber, Atoierr := strconv.Atoi(LDConfig[0])
					if Atoierr != nil {
						ctx.SendChain(message.Text("参与人数配置错误，本次抽奖无效，请通知机器人管理员"))
						var AllJoinUser string
						for _, i2 := range AllUserId {
							AllJoinUser += strconv.FormatInt(i2, 10) + ","
						}
						ctx.SendChain(message.Text("本次全部参与人员QQ号码记录\n" + AllJoinUser))
						return
					}
					LenAllUser := len(AllUserId)
					if LenAllUser < 1 {
						ctx.SendChain(message.Text("本次抽奖无人参与"))
						delete(LDLineUp, LDKeyString)
						return
					}
					for i := 0; i < LdNumber; i++ {
						ROD := NewPerm(0, LenAllUser)
						WinUser = append(WinUser, AllUserId[ROD])
						AllUserId = append(AllUserId[:ROD], AllUserId[ROD+1:]...)
					}
					ctx.SendChain(message.Text("关键词为 " + LDKeyString + " 的抽奖已结束，中奖用户为"))
					for _, i2 := range WinUser {
						ctx.SendChain(message.At(i2), message.Text("✨✨✨✨"), message.Text(i2), message.Image("http://q4.qlogo.cn/g?b=qq&nk="+strconv.FormatInt(i2, 10)+"&s=640").Add("cache", 0))
					}
					delete(LDLineUp, LDKeyString)
				}()
			}
		}
	})

	engine.OnPrefixGroup([]string{"参加抽奖", "参与抽奖"}, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		JLdText := ctx.State["args"].(string)
		JLdText = strings.Trim(JLdText, " ")
		_, IsEx := LDLineUp[JLdText]
		if !IsEx {
			ctx.SendChain(message.Text("关键词不存在，请检查关键词是否输入正确"))
		} else {

			UserId := ctx.Event.UserID
			for _, userId := range LDLineUp[JLdText] {
				if userId == UserId {
					ctx.SendChain(message.Text("您已参加同关键词的抽奖，请勿重复参加"))
					goto End
				}
			}

			LDLineUp[JLdText] = append(LDLineUp[JLdText], ctx.Event.UserID)
			fmt.Println(LDLineUp)
			ctx.SendChain(message.Text("已参与关键词为 "+JLdText+" 的抽奖"), message.At(ctx.Event.UserID))
		End:
		}
	})

}

func NewPerm(min, max int) int {
	//rand.Seed(time.Now().UnixNano())
	round := rand.Intn(max-min) + min
	return round
}
