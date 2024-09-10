package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ars0915/matching-system/constant"
	"github.com/ars0915/matching-system/entity"
	"github.com/ars0915/matching-system/util/cGin"
)

type addPersonBody struct {
	Name       string  `json:"name" binding:"required"`
	Height     float64 `json:"height" binding:"required"`
	Gender     string  `json:"gender" binding:"required"`
	WantedDate *uint64 `json:"wantedDate" binding:"required"`
}

func (rH *HttpHandler) addPersonAndFindMatchHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	var body addPersonBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid Json")
		return
	}

	data, err := rH.h.AddPersonAndFindMatch(ctx, entity.Person{
		Name:        body.Name,
		Height:      body.Height,
		Gender:      constant.Gender(body.Gender),
		WantedDates: body.WantedDate,
	})
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.WithData(data).Response(http.StatusOK, "")
}

func (rH *HttpHandler) printHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	rH.h.Print()

	ctx.Response(http.StatusOK, "")
}

func (rH *HttpHandler) removePersonHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid id")
		return
	}

	if err = rH.h.RemovePerson(ctx, uint64(id)); err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.Response(http.StatusOK, "")
}

func (rH *HttpHandler) querySinglePeopleHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid id")
		return
	}

	queryNumStr := ctx.Query("num")
	queryNum, err := strconv.Atoi(queryNumStr)
	if err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid query number")
		return
	}

	data, err := rH.h.QuerySinglePeople(ctx, uint64(id), queryNum)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.WithData(data).Response(http.StatusOK, "")
}
