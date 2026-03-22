package web

import (
	"errors"
	"io"
	"net/http"

	"qa_test_server/manager"

	"github.com/gin-gonic/gin"
)

type virtualStartRequest struct {
	Count       int    `json:"count"`
	IntervalMs  int    `json:"intervalMs"`
	Prefix      string `json:"prefix"`
	NamePrefix  string `json:"namePrefix"`
	StartIndex  int    `json:"startIndex"`
	Group       string `json:"group"`
	MutateParam bool   `json:"mutateParam"`
	WsBroadcast *bool  `json:"wsBroadcast"`
	PulseRepeat int    `json:"pulseRepeat"`
}

type virtualStopRequest struct {
	Remove *bool `json:"remove"`
}

func virtualStatus(c *gin.Context) {
	ok(c, manager.VirtualDeviceManagerGlobal.Status())
}

func virtualStart(c *gin.Context) {
	req := virtualStartRequest{}
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		fail(c, http.StatusBadRequest, "invalid start payload")
		return
	}

	wsBroadcast := true
	if req.WsBroadcast != nil {
		wsBroadcast = *req.WsBroadcast
	}

	status, err := manager.VirtualDeviceManagerGlobal.Start(manager.VirtualDeviceConfig{
		Count:       req.Count,
		IntervalMs:  req.IntervalMs,
		Prefix:      req.Prefix,
		NamePrefix:  req.NamePrefix,
		StartIndex:  req.StartIndex,
		Group:       req.Group,
		MutateParam: req.MutateParam,
		WsBroadcast: wsBroadcast,
	})
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	ok(c, status)
}

func virtualStop(c *gin.Context) {
	req := virtualStopRequest{}
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		fail(c, http.StatusBadRequest, "invalid stop payload")
		return
	}
	remove := true
	if req.Remove != nil {
		remove = *req.Remove
	}
	status, err := manager.VirtualDeviceManagerGlobal.Stop(remove)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, status)
}

func virtualPulse(c *gin.Context) {
	req := virtualStartRequest{}
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		fail(c, http.StatusBadRequest, "invalid pulse payload")
		return
	}

	wsBroadcast := true
	if req.WsBroadcast != nil {
		wsBroadcast = *req.WsBroadcast
	}
	result, err := manager.VirtualDeviceManagerGlobal.Pulse(manager.VirtualDeviceConfig{
		Count:       req.Count,
		IntervalMs:  req.IntervalMs,
		Prefix:      req.Prefix,
		NamePrefix:  req.NamePrefix,
		StartIndex:  req.StartIndex,
		Group:       req.Group,
		MutateParam: req.MutateParam,
		WsBroadcast: wsBroadcast,
		PulseRepeat: req.PulseRepeat,
	})
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	ok(c, result)
}

func virtualStressPulse(c *gin.Context) {
	req := virtualStartRequest{}
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		fail(c, http.StatusBadRequest, "invalid stress payload")
		return
	}

	wsBroadcast := true
	if req.WsBroadcast != nil {
		wsBroadcast = *req.WsBroadcast
	}

	result, err := manager.VirtualDeviceManagerGlobal.StressPulse(manager.VirtualDeviceConfig{
		Count:       req.Count,
		IntervalMs:  req.IntervalMs,
		Prefix:      req.Prefix,
		NamePrefix:  req.NamePrefix,
		StartIndex:  req.StartIndex,
		Group:       req.Group,
		MutateParam: req.MutateParam,
		WsBroadcast: wsBroadcast,
		PulseRepeat: req.PulseRepeat,
	})
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	ok(c, result)
}
