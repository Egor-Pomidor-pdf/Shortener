package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/dto"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/models"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/service"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(srv service.ServiceInterface) *Handler {
	return &Handler{service: srv}
}



func (h *Handler) postShorten(c *ginext.Context) { 
	var body dto.ShortenRequest
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ginext.H{"error": fmt.Sprintf("invalid body (parsing): %s", err.Error())})
		return
	}

	var createLink *models.ShortURL
	createLink, err = body.ToEntity()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ginext.H{"error": fmt.Sprintf("invalid body (validating): %s", err.Error())})
		return
	}

	shortLink, err := h.service.CreateShort(context.Background(), createLink)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusConflict,
			ginext.H{"error": fmt.Sprintf("couldn't perform operation: %s", err.Error())},
		)
		return
	}
	

	c.JSON(http.StatusCreated, shortLink)
}


func (h *Handler) getRedirect(c *ginext.Context) {
    short := c.Param("short")

    s, err := h.service.Resolve(c.Request.Context(), short)
    if err != nil || s == nil {
        c.AbortWithStatus(http.StatusNotFound)
        return
    }

    var clientID *uuid.UUID
    if cid := c.GetHeader("X-Client-Id"); cid != "" {
        if u, err := uuid.Parse(cid); err == nil {
            clientID = &u
        }
    }
    if clientID == nil {
        if cid := c.Query("client_id"); cid != "" {
            if u, err := uuid.Parse(cid); err == nil {
                clientID = &u
            }
        }
    }

    if clientID != nil {
        ip := c.ClientIP()

        _ = h.service.RecordClick(c.Request.Context(), models.ClickEvent{
            ShortCode: short,
            ClientID:  *clientID,
            UserAgent: c.Request.UserAgent(),
            IP:        ip,
            At:        time.Now(),
        })
    }

    c.Redirect(http.StatusFound, s.Original)
}

type analyticsResponse struct {
    Total uint64             `json:"total"`
    Daily []models.AggPoint `json:"daily"`
    ByUA  []models.AggPoint `json:"by_user_agent"`
}

func (h *Handler) getAnalytics(c *ginext.Context) {
    short := c.Param("short")
    ctx := c.Request.Context()

    total, err := h.service.Count(ctx, short)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, ginext.H{
            "error": err.Error(),
        })
        return
    }

    daily, err := h.service.Daily(ctx, short)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, ginext.H{
            "error": err.Error(),
        })
        return
    }

    byUA, err := h.service.ByUserAgent(ctx, short)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, ginext.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, analyticsResponse{
        Total: total,
        Daily: daily,
        ByUA:  byUA,
    })
}
