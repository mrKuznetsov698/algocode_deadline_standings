package server

import (
	processors "algocode_deadline_standings/dataProcessor"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

//func getData(config *configs.Config, needStats bool) (
//	[]*processors.CriterionTitle,
//	[]*processors.UserValues,
//	map[int]*processors.Stats,
//	error) {
//
//	// hope we can do it on every request
//	criterionTitles, userValues, err := processors.GetDeadlineResults(config)
//	var stats map[int]*processors.Stats
//	if needStats && err == nil {
//		stats, err = processors.CreateStatistics(config, userValues)
//	}
//	// yes, that will return trash in stats if we don't need it or an error occurred
//	return criterionTitles, userValues, stats, err
//}

//func dataTick(c *gin.Context, data *DataAccess) {
//	if err := data.tick(); err != nil {
//		c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
//			"Error": err.Error(),
//		})
//		c.Abort()
//	}
//}

func addHandler(handler func(*gin.Context, *DataAccess), data *DataAccess) GinHandler {
	return func(c *gin.Context) {
		if err := data.tick(); err != nil {
			c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
				"Error": err.Error(),
			})
			return
		}
		data.lock.RLock()
		defer data.lock.RUnlock()
		handler(c, data)
	}
}

func mainPage(c *gin.Context, data *DataAccess) {
	c.HTML(http.StatusOK, "page.gohtml", gin.H{
		"CriterionTitles": data.criterionTitles,
		"UserValues":      data.userValues,
		"Single":          false,
	})
}

func studentStats(c *gin.Context, data *DataAccess) {
	name := c.Param("name")
	ind, found := slices.BinarySearchFunc(data.userValues, name, func(values *processors.UserValues, s string) int {
		return strings.Compare(values.FullName, s)
	})
	if found {
		c.HTML(http.StatusOK, "page.gohtml", gin.H{
			"CriterionTitles": data.criterionTitles,
			"UserValues":      []*processors.UserValues{data.userValues[ind]},
			"Single":          true,
		})
	} else {
		c.String(http.StatusNotFound, "Nothing found with name=\"%s\"", name)
	}
}

func allStats(c *gin.Context, data *DataAccess) {
	c.HTML(http.StatusOK, "stats.gohtml", gin.H{
		"Stats": data.stats,
	})
}
