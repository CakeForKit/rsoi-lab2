package controller

import (
	"sync"

	coreController "github.com/CakeForKit/rsoi-lab2/lb-core/controller"
	"github.com/CakeForKit/rsoi-lab2/lb-core/model"
	"github.com/CakeForKit/rsoi-lab2/ms-gateway/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var apiCntrlr coreController.HttpController
var apiCntrlrMutex sync.Mutex

type apiController struct {
	service *service.GatewayService
}

func GetAPIController() (coreController.HttpController, error) {
	apiCntrlrMutex.Lock()
	defer apiCntrlrMutex.Unlock()

	if apiCntrlr != nil {
		return apiCntrlr, nil
	}

	apiCntrlr = &apiController{service: service.New()}
	return apiCntrlr, nil
}

func (a *apiController) RegisterHttpController(r *gin.Engine) {
	g := r.Group("/api/v1")
	g.GET("/flights", a.flights)
	g.GET("/privilege", a.privilege)
	g.POST("/tickets", a.purchase)
	g.GET("/tickets", a.tickets)
	g.GET("/tickets/:ticketUid", a.ticket)
	g.DELETE("/tickets/:ticketUid", a.cancel)
	g.GET("/me", a.me)
}
func username(c *gin.Context) (string, bool) {
	u := c.GetHeader("X-User-Name")
	if u == "" {
		c.JSON(400, gin.H{"message": "X-User-Name is required"})
		return "", false
	}
	return u, true
}
func failure(c *gin.Context, e error) { c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()}) }
func (a *apiController) flights(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	s, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	v, e := a.service.Flights(p, s)
	if e != nil {
		failure(c, e)
		return
	}
	c.JSON(200, v)
}
func (a *apiController) privilege(c *gin.Context) {
	u, ok := username(c)
	if !ok {
		return
	}
	v, e := a.service.GetBonus(u)
	if e != nil {
		failure(c, e)
		return
	}
	c.JSON(200, v)
}
func (a *apiController) purchase(c *gin.Context) {
	u, ok := username(c)
	if !ok {
		return
	}
	var req model.TicketPurchaseRequest
	if c.ShouldBindJSON(&req) != nil {
		failure(c, http.ErrNotSupported)
		return
	}
	v, e := a.service.Purchase(u, req)
	if e != nil {
		failure(c, e)
		return
	}
	c.JSON(200, v)
}
func (a *apiController) tickets(c *gin.Context) {
	u, ok := username(c)
	if !ok {
		return
	}
	v, e := a.service.Tickets(u)
	if e != nil {
		failure(c, e)
		return
	}
	c.JSON(200, v)
}
func (a *apiController) ticket(c *gin.Context) {
	u, ok := username(c)
	if !ok {
		return
	}
	v, e := a.service.Ticket(u, c.Param("ticketUid"))
	if e != nil {
		c.JSON(404, gin.H{"message": "ticket not found"})
		return
	}
	c.JSON(200, v)
}
func (a *apiController) cancel(c *gin.Context) {
	u, ok := username(c)
	if !ok {
		return
	}
	if e := a.service.Cancel(u, c.Param("ticketUid")); e != nil {
		c.JSON(404, gin.H{"message": e.Error()})
		return
	}
	c.Status(204)
}
func (a *apiController) me(c *gin.Context) {
	u, ok := username(c)
	if !ok {
		return
	}
	v, e := a.service.Me(u)
	if e != nil {
		failure(c, e)
		return
	}
	c.JSON(200, v)
}
