package resultshandling

import (
	"encoding/json"

	"github.com/armosec/kubescape/v2/core/cautils"
	"github.com/armosec/kubescape/v2/core/cautils/logger"
	"github.com/armosec/kubescape/v2/core/cautils/logger/helpers"
	"github.com/armosec/kubescape/v2/core/pkg/resultshandling/printer"
	printerv1 "github.com/armosec/kubescape/v2/core/pkg/resultshandling/printer/v1"
	printerv2 "github.com/armosec/kubescape/v2/core/pkg/resultshandling/printer/v2"
	"github.com/armosec/kubescape/v2/core/pkg/resultshandling/reporter"
)

type ResultsHandler struct {
	reporterObj reporter.IReport
	printerObj  printer.IPrinter
	scanData    *cautils.OPASessionObj
}

func NewResultsHandler(reporterObj reporter.IReport, printerObj printer.IPrinter) *ResultsHandler {
	return &ResultsHandler{
		reporterObj: reporterObj,
		printerObj:  printerObj,
	}
}

// GetScore return scan risk-score
func (resultsHandler *ResultsHandler) GetRiskScore() float32 {
	return resultsHandler.scanData.Report.SummaryDetails.Score
}

// GetData get scan/action related data (policies, resources, results, etc.). Call ToJson function if you wish the json representation of the data
func (resultsHandler *ResultsHandler) GetData() *cautils.OPASessionObj {
	return resultsHandler.scanData
}

// SetData set scan/action related data
func (resultsHandler *ResultsHandler) SetData(data *cautils.OPASessionObj) {
	resultsHandler.scanData = data
}

// GetPrinter get printer object
func (resultsHandler *ResultsHandler) GetPrinter() printer.IPrinter {
	return resultsHandler.printerObj
}

// GetReporter get reporter object
func (resultsHandler *ResultsHandler) GetReporter() reporter.IReport {
	return resultsHandler.reporterObj
}

// ToJson return results in json format
func (resultsHandler *ResultsHandler) ToJson() ([]byte, error) {
	return json.Marshal(printerv2.DataToJson(resultsHandler.scanData))
}

// HandleResults handle the scan results according to the pre defind interfaces
func (resultsHandler *ResultsHandler) HandleResults() error {

	resultsHandler.printerObj.ActionPrint(resultsHandler.scanData)

	if err := resultsHandler.reporterObj.Submit(resultsHandler.scanData); err != nil {
		return err
	}

	resultsHandler.printerObj.Score(resultsHandler.GetRiskScore())

	resultsHandler.reporterObj.DisplayReportURL()

	return nil
}

// NewPrinter defind output format
func NewPrinter(printFormat, formatVersion string, verboseMode bool, viewType cautils.ViewTypes) printer.IPrinter {

	switch printFormat {
	case printer.JsonFormat:
		switch formatVersion {
		case "v2":
			return printerv2.NewJsonPrinter()
		default:
			logger.L().Warning("Deprecated format version", helpers.String("run", "--format-version=v2"))
			return printerv1.NewJsonPrinter()
		}
	case printer.JunitResultFormat:
		return printerv2.NewJunitPrinter(verboseMode)
	case printer.PrometheusFormat:
		return printerv2.NewPrometheusPrinter(verboseMode)
	case printer.PdfFormat:
		return printerv2.NewPdfPrinter()
	default:
		return printerv2.NewPrettyPrinter(verboseMode, formatVersion, viewType)
	}
}
