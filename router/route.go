package router

//router.go : api 전체 인입에 대한 관리 및 구성을 담당하는 파일
import (
	"fmt"
	ctl "lecture/go-contracts/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ct *ctl.Controller
}

func NewRouter(ctl *ctl.Controller) (*Router, error) {
	r := &Router{ct: ctl} //controller 포인터를 ct로 복사, 할당
	return r, nil
}

// cross domain을 위해 사용
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//허용할 header 타입에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		//허용할 method에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// 임의 인증을 위한 함수
func liteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c == nil {
			c.Abort() // 미들웨어에서 사용, 이후 요청 중지
			return
		}
		//http 헤더내 "Authorization" 폼의 데이터를 조회
		auth := c.GetHeader("Authorization")
		//실제 인증기능이 올수있다. 단순히 출력기능만 처리 현재는 출력예시
		fmt.Println("Authorization-word ", auth)

		c.Next() // 다음 요청 진행
	}
}

// 실제 라우팅
func (p *Router) Idx() *gin.Engine {
	e := gin.New() // gin선언

	//GET 라우팅
	view_contract := e.Group("/view", liteAuth())
	{
		view_contract.GET("/getSymbol/:tokenName", p.ct.GetTokenSymbol)          // 심볼 조회
		view_contract.GET("/getBalance/:address", p.ct.GetTokenBalanceByAddress) // 자산 조회
	}
	//POST  라우팅
	post_contract := e.Group("/post", liteAuth())
	{
		post_contract.POST("/coinTransfer/:address", p.ct.CoinTransfer)           //코인 전송
		post_contract.POST("/coinTransferFrom/:address", p.ct.CoinTransferFrom)   //다른 개인 키로, 특정 주소에 코인 전송
		post_contract.POST("/tokenTransfer/:address", p.ct.TokenTransfer)         //토큰 전송
		post_contract.POST("/tokenTransferFrom/:address", p.ct.TokenTransferFrom) //다른 개인 키로, 특정 주소에 토큰 전송
	}

	return e
}
