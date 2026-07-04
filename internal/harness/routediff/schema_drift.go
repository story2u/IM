package routediff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"wework-go/internal/contracts"
	"wework-go/internal/httpserver"
	"wework-go/internal/inventory"
)

const maxSchemaDriftRows = 60
const maxSchemaDriftReasonLines = 12

// SchemaDriftReasonCount aggregates drift reasons for reporting.
type SchemaDriftReasonCount struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

// SchemaDriftRow captures schema-level compatibility status for one route.
type SchemaDriftRow struct {
	Method                    string   `json:"method"`
	Path                      string   `json:"path"`
	Owner                     string   `json:"owner"`
	Phase                     string   `json:"phase"`
	ReferenceResponseContract string   `json:"reference_response_contract"`
	GoResponseContract        string   `json:"go_response_contract"`
	ReferenceRequestContract  string   `json:"reference_request_contract"`
	GoRequestContract         string   `json:"go_request_contract"`
	ReferenceResponseSig      string   `json:"reference_response_schema_signature"`
	GoResponseSig             string   `json:"go_response_schema_signature"`
	ReferenceRequestSig       string   `json:"reference_request_schema_signature"`
	GoRequestSig              string   `json:"go_request_schema_signature"`
	ResponseSchemaMatch       bool     `json:"response_schema_match"`
	ResponseSchemaReasons     []string `json:"response_schema_reasons"`
	RequestSchemaMatch        bool     `json:"request_schema_match"`
	RequestSchemaReasons      []string `json:"request_schema_reasons"`
	ResponseSchemaDiff        []string `json:"response_schema_diff"`
	RequestSchemaDiff         []string `json:"request_schema_diff"`
}

// SchemaDriftReport summarizes contract-level route parity.
type SchemaDriftReport struct {
	ReferenceRouteCount          int                      `json:"reference_route_count"`
	GoRouteCount                 int                      `json:"go_route_count"`
	MatchingCount                int                      `json:"matching_count"`
	ReferenceOnlyCount           int                      `json:"reference_only_count"`
	GoOnlyCount                  int                      `json:"go_only_count"`
	SchemaComparableCount        int                      `json:"schema_comparable_count"`
	SchemaMatchCount             int                      `json:"schema_match_count"`
	SchemaMismatchCount          int                      `json:"schema_mismatch_count"`
	MissingReferenceContractLink int                      `json:"missing_reference_contract_link_count"`
	MissingGoContractLink        int                      `json:"missing_go_contract_link_count"`
	TopDriftReasons              []SchemaDriftReasonCount `json:"top_drift_reasons"`
	Rows                         []SchemaDriftRow         `json:"rows"`
}

// BuildSchemaDriftReport builds a contract-shape report from reference and go routes.
func BuildSchemaDriftReport(referenceRoutes []inventory.Route, goRoutes []httpserver.Route, catalog []contracts.SchemaFile) SchemaDriftReport {
	report := CompareWithContracts(referenceRoutes, goRoutes, catalog)
	contractIndex := buildSchemaContractIndex(catalog)
	signatures := make(map[string]string, len(contractIndex))

	result := SchemaDriftReport{
		ReferenceRouteCount: report.ReferenceRouteCount,
		GoRouteCount:        report.GoRouteCount,
		MatchingCount:       len(report.Matching),
		ReferenceOnlyCount:  len(report.ReferenceOnly),
		GoOnlyCount:         len(report.GoOnly),
		TopDriftReasons:     nil,
	}

	reasonCounter := make(map[string]int, 16)
	for _, row := range report.Matching {
		maybeEmitRow := false
		driftRow := buildSchemaDriftRow(row, contractIndex, signatures, reasonCounter)
		if len(row.ReferenceResponseContract) > 0 || len(row.GoResponseContract) > 0 || len(row.ReferenceRequestContract) > 0 || len(row.GoRequestContract) > 0 {
			result.SchemaComparableCount++
			if row.ReferenceResponseContract == "" {
				result.MissingReferenceContractLink++
			}
			if row.GoResponseContract == "" {
				result.MissingGoContractLink++
			}
			if row.ReferenceRequestContract == "" {
				result.MissingReferenceContractLink++
			}
			if row.GoRequestContract == "" {
				result.MissingGoContractLink++
			}
			maybeEmitRow = true
		}
		if maybeEmitRow && len(driftRow.ResponseSchemaReasons) == 0 && len(driftRow.RequestSchemaReasons) == 0 {
			result.SchemaMatchCount++
			continue
		}
		if maybeEmitRow {
			result.SchemaMismatchCount++
			result.Rows = append(result.Rows, driftRow)
			if len(result.Rows) > maxSchemaDriftRows {
				result.Rows = result.Rows[:maxSchemaDriftRows]
			}
		}
	}

	for _, reason := range sortedSchemaDriftReasons(reasonCounter) {
		result.TopDriftReasons = append(result.TopDriftReasons, SchemaDriftReasonCount{
			Reason: reason.reason,
			Count:  reason.count,
		})
	}
	return result
}

func buildSchemaDriftRow(routeRef RouteRef, contractIndex map[string]contracts.SchemaFile, signatures map[string]string, reasonCounter map[string]int) SchemaDriftRow {
	row := SchemaDriftRow{
		Method:                    routeRef.Method,
		Path:                      routeRef.Path,
		Owner:                     routeRef.Owner,
		Phase:                     routeRef.Phase,
		ReferenceResponseContract: routeRef.ReferenceResponseContract,
		GoResponseContract:        routeRef.GoResponseContract,
		ReferenceRequestContract:  routeRef.ReferenceRequestContract,
		GoRequestContract:         routeRef.GoRequestContract,
		ReferenceResponseSig:      schemaSignature(routeRef.ReferenceResponseContract, contractIndex, signatures),
		GoResponseSig:             schemaSignature(routeRef.GoResponseContract, contractIndex, signatures),
		ReferenceRequestSig:       schemaSignature(routeRef.ReferenceRequestContract, contractIndex, signatures),
		GoRequestSig:              schemaSignature(routeRef.GoRequestContract, contractIndex, signatures),
	}

	if !hasAnyContract(row.ReferenceResponseContract, row.GoResponseContract, row.ReferenceRequestContract, row.GoRequestContract) {
		return row
	}

	row.ResponseSchemaMatch = true
	row.RequestSchemaMatch = true

	if row.ReferenceResponseContract != "" || row.GoResponseContract != "" {
		ok, reasons, diffs := compareSchemaContracts(routeRef.ReferenceResponseContract, routeRef.GoResponseContract, "response", contractIndex, signatures)
		row.ResponseSchemaMatch = ok
		row.ResponseSchemaReasons = reasons
		row.ResponseSchemaDiff = diffs
		countReasons(reasonCounter, reasons)
	}
	if row.ReferenceRequestContract != "" || row.GoRequestContract != "" {
		ok, reasons, diffs := compareSchemaContracts(routeRef.ReferenceRequestContract, routeRef.GoRequestContract, "request", contractIndex, signatures)
		row.RequestSchemaMatch = ok
		row.RequestSchemaReasons = reasons
		row.RequestSchemaDiff = diffs
		countReasons(reasonCounter, reasons)
	}

	row.ResponseSchemaMatch = len(row.ResponseSchemaReasons) == 0
	row.RequestSchemaMatch = len(row.RequestSchemaReasons) == 0
	return row
}

type schemaDriftReasonAgg struct {
	reason string
	count  int
}

func compareSchemaContracts(referenceName, goName, kind string, contractIndex map[string]contracts.SchemaFile, signatures map[string]string) (bool, []string, []string) {
	reasons := []string{}
	diffs := []string{}

	if referenceName == "" && goName == "" {
		return true, reasons, diffs
	}
	if referenceName == "" {
		return false, []string{"reference " + kind + " contract missing"}, diffs
	}
	if goName == "" {
		return false, []string{"go " + kind + " contract missing"}, diffs
	}
	if referenceName != goName {
		reasons = append(reasons, "contract name mismatch")
	}

	referenceKey := canonicalSchemaKey(referenceName)
	goKey := canonicalSchemaKey(goName)
	left, leftExists := contractIndex[referenceKey]
	right, rightExists := contractIndex[goKey]
	if !leftExists || strings.TrimSpace(left.Path) == "" {
		reasons = append(reasons, "reference "+kind+" contract definition missing")
	}
	if !rightExists || strings.TrimSpace(right.Path) == "" {
		reasons = append(reasons, "go "+kind+" contract definition missing")
	}
	if len(reasons) > 0 {
		return false, dedupeSortedReasons(reasons), diffs
	}

	leftSig := schemaSignature(referenceName, contractIndex, signatures)
	rightSig := schemaSignature(goName, contractIndex, signatures)
	if leftSig != rightSig {
		diffs = schemaDiff(left.Path, right.Path)
		reasons = append(reasons, "schema shape mismatch")
		return false, dedupeSortedReasons(reasons), dedupeTrimDiffs(diffs)
	}
	return true, nil, nil
}

func schemaDiff(leftPath, rightPath string) []string {
	leftBytes, err := os.ReadFile(filepath.Clean(leftPath))
	if err != nil {
		return []string{fmt.Sprintf("read schema failed: %v", err)}
	}
	rightBytes, err := os.ReadFile(filepath.Clean(rightPath))
	if err != nil {
		return []string{fmt.Sprintf("read schema failed: %v", err)}
	}

	var left, right any
	if err := json.Unmarshal(leftBytes, &left); err != nil {
		return []string{fmt.Sprintf("decode schema %q failed: %v", filepath.Base(leftPath), err)}
	}
	if err := json.Unmarshal(rightBytes, &right); err != nil {
		return []string{fmt.Sprintf("decode schema %q failed: %v", filepath.Base(rightPath), err)}
	}
	left = normalizeSchemaNode(left)
	right = normalizeSchemaNode(right)

	diffs := []string{}
	collectSchemaDiffs("$", left, right, &diffs)
	return dedupeTrimDiffs(diffs)
}

func collectSchemaDiffs(path string, left any, right any, diffs *[]string) {
	if len(*diffs) >= maxSchemaDriftReasonLines {
		return
	}
	if reflect.DeepEqual(left, right) {
		return
	}

	leftObj := left
	rightObj := right
	leftMap, leftIsMap := leftObj.(map[string]any)
	rightMap, rightIsMap := rightObj.(map[string]any)
	if leftIsMap && rightIsMap {
		keys := map[string]struct{}{}
		for key := range leftMap {
			if shouldSkipSchemaKey(key) {
				continue
			}
			keys[key] = struct{}{}
		}
		for key := range rightMap {
			if shouldSkipSchemaKey(key) {
				continue
			}
			keys[key] = struct{}{}
		}
		allKeys := make([]string, 0, len(keys))
		for key := range keys {
			allKeys = append(allKeys, key)
		}
		sort.Strings(allKeys)

		for _, key := range allKeys {
			if len(*diffs) >= maxSchemaDriftReasonLines {
				return
			}
			childPath := path + "." + key
			leftChild, leftHas := leftMap[key]
			rightChild, rightHas := rightMap[key]
			if !leftHas {
				*diffs = append(*diffs, fmt.Sprintf("%s: missing in reference contract", childPath))
				continue
			}
			if !rightHas {
				*diffs = append(*diffs, fmt.Sprintf("%s: missing in go contract", childPath))
				continue
			}
			collectSchemaDiffs(childPath, leftChild, rightChild, diffs)
		}
		return
	}

	leftArr, leftIsArray := left.([]any)
	rightArr, rightIsArray := right.([]any)
	if leftIsArray || rightIsArray {
		if !leftIsArray || !rightIsArray {
			*diffs = append(*diffs, fmt.Sprintf("%s: type changed from %T to %T", path, leftObj, rightObj))
			return
		}
		if len(leftArr) != len(rightArr) {
			*diffs = append(*diffs, fmt.Sprintf("%s: array length %d vs %d", path, len(leftArr), len(rightArr)))
		}
		max := len(leftArr)
		if len(rightArr) > max {
			max = len(rightArr)
		}
		for idx := 0; idx < max && len(*diffs) < maxSchemaDriftReasonLines; idx++ {
			childPath := fmt.Sprintf("%s[%d]", path, idx)
			if idx >= len(leftArr) {
				*diffs = append(*diffs, fmt.Sprintf("%s: extra value in go contract: %v", childPath, rightArr[idx]))
				continue
			}
			if idx >= len(rightArr) {
				*diffs = append(*diffs, fmt.Sprintf("%s: extra value in reference contract: %v", childPath, leftArr[idx]))
				continue
			}
			collectSchemaDiffs(childPath, leftArr[idx], rightArr[idx], diffs)
		}
		return
	}

	if reflect.TypeOf(leftObj) != reflect.TypeOf(rightObj) {
		*diffs = append(*diffs, fmt.Sprintf("%s: type changed from %T to %T", path, leftObj, rightObj))
		return
	}

	if !reflect.DeepEqual(leftObj, rightObj) {
		*diffs = append(*diffs, fmt.Sprintf("%s: %v vs %v", path, leftObj, rightObj))
	}
}

func dedupeSortedReasons(reasons []string) []string {
	uniq := make(map[string]struct{}, len(reasons))
	out := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		if reason == "" {
			continue
		}
		if _, ok := uniq[reason]; ok {
			continue
		}
		uniq[reason] = struct{}{}
		out = append(out, reason)
	}
	sort.Strings(out)
	return out
}

func dedupeTrimDiffs(diffs []string) []string {
	uniq := make(map[string]struct{}, len(diffs))
	out := make([]string, 0, len(diffs))
	for _, diff := range diffs {
		diff = strings.TrimSpace(diff)
		if diff == "" {
			continue
		}
		if _, ok := uniq[diff]; ok {
			continue
		}
		uniq[diff] = struct{}{}
		out = append(out, diff)
		if len(out) >= maxSchemaDriftReasonLines {
			break
		}
	}
	sort.Strings(out)
	return out
}

func hasAnyContract(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}

func containsReason(reasons []string, target string) bool {
	for _, reason := range reasons {
		if reason == target {
			return true
		}
	}
	return false
}

func countReasons(counter map[string]int, reasons []string) {
	for _, reason := range reasons {
		if reason == "" {
			continue
		}
		counter[reason]++
	}
}

func sortedSchemaDriftReasons(counter map[string]int) []schemaDriftReasonAgg {
	list := make([]schemaDriftReasonAgg, 0, len(counter))
	for reason, count := range counter {
		list = append(list, schemaDriftReasonAgg{reason: reason, count: count})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].count == list[j].count {
			return list[i].reason < list[j].reason
		}
		return list[i].count > list[j].count
	})
	return list
}
