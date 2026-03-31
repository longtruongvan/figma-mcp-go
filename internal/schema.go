package internal

import (
	"fmt"
	"regexp"
)

// nodeIDPattern matches Figma node IDs: colon-separated integers e.g. "4029:12345"
var nodeIDPattern = regexp.MustCompile(`^\d+:\d+$`)

// ValidNodeID reports whether s is a valid Figma node ID.
func ValidNodeID(s string) bool {
	return nodeIDPattern.MatchString(s)
}

// ValidateRPC validates an incoming RPC request against the tool's expected
// input shape. Returns an error string on failure, empty string if valid.
func ValidateRPC(tool string, nodeIDs []string, params map[string]interface{}) string {
	switch tool {
	case "get_node":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}

	case "get_nodes_info":
		if len(nodeIDs) == 0 {
			return "nodeIds is required and must not be empty"
		}
		for _, id := range nodeIDs {
			if !ValidNodeID(id) {
				return fmt.Sprintf("invalid nodeId: %s — must use colon format e.g. 4029:12345", id)
			}
		}

	case "get_screenshot":
		for _, id := range nodeIDs {
			if !ValidNodeID(id) {
				return fmt.Sprintf("invalid nodeId: %s — must use colon format e.g. 4029:12345", id)
			}
		}
		if format, ok := params["format"].(string); ok {
			if !validExportFormat(format) {
				return fmt.Sprintf("format must be PNG, SVG, JPG, or PDF, got: %s", format)
			}
		}

	case "save_screenshots":
		items, ok := params["items"]
		if !ok {
			return "items is required"
		}
		itemList, ok := items.([]interface{})
		if !ok || len(itemList) == 0 {
			return "items must be a non-empty array"
		}
		for i, item := range itemList {
			m, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Sprintf("items[%d] must be an object", i)
			}
			nodeID, _ := m["nodeId"].(string)
			if !ValidNodeID(nodeID) {
				return fmt.Sprintf("items[%d].nodeId must use colon format e.g. 4029:12345", i)
			}
			outputPath, _ := m["outputPath"].(string)
			if outputPath == "" {
				return fmt.Sprintf("items[%d].outputPath is required", i)
			}
		}

	case "get_design_context":
		if depth, ok := params["depth"].(float64); ok {
			if depth < 0 {
				return "depth must be a non-negative number"
			}
		}
		if detail, ok := params["detail"].(string); ok && detail != "" {
			switch detail {
			case "minimal", "compact", "full":
			default:
				return fmt.Sprintf("detail must be minimal, compact, or full, got: %s", detail)
			}
		}

	case "search_nodes":
		query, _ := params["query"].(string)
		if query == "" {
			return "query is required"
		}
		if nodeID, ok := params["nodeId"].(string); ok && nodeID != "" {
			if !ValidNodeID(nodeID) {
				return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeID)
			}
		}
		if limit, ok := params["limit"].(float64); ok && limit <= 0 {
			return "limit must be a positive number"
		}

	case "get_reactions":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}

	case "scan_text_nodes", "scan_nodes_by_types":
		nodeID, _ := params["nodeId"].(string)
		if nodeID == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeID) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeID)
		}
		if tool == "scan_nodes_by_types" {
			types, ok := params["types"].([]interface{})
			if !ok || len(types) == 0 {
				return "types must be a non-empty array"
			}
		}

	// ── Write tools ──────────────────────────────────────────────────────────

	case "create_frame":
		if w, ok := params["width"].(float64); ok && w <= 0 {
			return "width must be positive"
		}
		if h, ok := params["height"].(float64); ok && h <= 0 {
			return "height must be positive"
		}
		if lm, ok := params["layoutMode"].(string); ok && lm != "" {
			switch lm {
			case "HORIZONTAL", "VERTICAL", "NONE":
			default:
				return fmt.Sprintf("layoutMode must be HORIZONTAL, VERTICAL, or NONE, got: %s", lm)
			}
		}
		if pid, ok := params["parentId"].(string); ok && pid != "" && !ValidNodeID(pid) {
			return fmt.Sprintf("parentId must use colon format e.g. 4029:12345, got: %s", pid)
		}

	case "create_rectangle", "create_ellipse":
		if w, ok := params["width"].(float64); ok && w <= 0 {
			return "width must be positive"
		}
		if h, ok := params["height"].(float64); ok && h <= 0 {
			return "height must be positive"
		}
		if pid, ok := params["parentId"].(string); ok && pid != "" && !ValidNodeID(pid) {
			return fmt.Sprintf("parentId must use colon format e.g. 4029:12345, got: %s", pid)
		}

	case "create_text":
		if text, _ := params["text"].(string); text == "" {
			return "text is required"
		}
		if v, ok := params["fontSize"].(float64); ok && v <= 0 {
			return "fontSize must be positive"
		}
		if v, ok := params["width"].(float64); ok && v <= 0 {
			return "width must be positive"
		}
		if v, ok := params["height"].(float64); ok && v <= 0 {
			return "height must be positive"
		}
		if err := validateTextAutoResize(params["textAutoResize"]); err != "" {
			return err
		}
		if v, ok := params["textAlignHorizontal"].(string); ok && v != "" {
			switch v {
			case "LEFT", "CENTER", "RIGHT", "JUSTIFIED":
			default:
				return fmt.Sprintf("textAlignHorizontal must be LEFT, CENTER, RIGHT, or JUSTIFIED, got: %s", v)
			}
		}
		if v, ok := params["textAlignVertical"].(string); ok && v != "" {
			switch v {
			case "TOP", "CENTER", "BOTTOM":
			default:
				return fmt.Sprintf("textAlignVertical must be TOP, CENTER, or BOTTOM, got: %s", v)
			}
		}
		if v, ok := params["paragraphSpacing"].(float64); ok && v < 0 {
			return "paragraphSpacing must be non-negative"
		}
		if err := validateLineHeight(params["lineHeight"]); err != "" {
			return err
		}
		if err := validateLetterSpacing(params["letterSpacing"]); err != "" {
			return err
		}
		if pid, ok := params["parentId"].(string); ok && pid != "" && !ValidNodeID(pid) {
			return fmt.Sprintf("parentId must use colon format e.g. 4029:12345, got: %s", pid)
		}

	case "set_text":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if _, ok := params["text"].(string); !ok {
			return "text is required"
		}

	case "set_fills":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if color, _ := params["color"].(string); color == "" {
			return "color is required (hex string e.g. #FF5733)"
		}

	case "set_strokes":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if color, _ := params["color"].(string); color == "" {
			return "color is required (hex string e.g. #FF5733)"
		}

	case "move_nodes":
		if len(nodeIDs) == 0 {
			return "nodeIds is required"
		}
		for _, id := range nodeIDs {
			if !ValidNodeID(id) {
				return fmt.Sprintf("invalid nodeId: %s — must use colon format e.g. 4029:12345", id)
			}
		}
		_, hasX := params["x"]
		_, hasY := params["y"]
		if !hasX && !hasY {
			return "at least one of x or y is required"
		}

	case "resize_nodes":
		if len(nodeIDs) == 0 {
			return "nodeIds is required"
		}
		for _, id := range nodeIDs {
			if !ValidNodeID(id) {
				return fmt.Sprintf("invalid nodeId: %s — must use colon format e.g. 4029:12345", id)
			}
		}
		_, hasW := params["width"]
		_, hasH := params["height"]
		if !hasW && !hasH {
			return "at least one of width or height is required"
		}

	case "delete_nodes":
		if len(nodeIDs) == 0 {
			return "nodeIds is required and must not be empty"
		}
		for _, id := range nodeIDs {
			if !ValidNodeID(id) {
				return fmt.Sprintf("invalid nodeId: %s — must use colon format e.g. 4029:12345", id)
			}
		}

	case "rename_node":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if name, _ := params["name"].(string); name == "" {
			return "name is required"
		}

	case "clone_node":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if pid, ok := params["parentId"].(string); ok && pid != "" && !ValidNodeID(pid) {
			return fmt.Sprintf("parentId must use colon format e.g. 4029:12345, got: %s", pid)
		}

	case "import_image":
		if imageData, _ := params["imageData"].(string); imageData == "" {
			return "imageData (base64) is required"
		}
		if sm, ok := params["scaleMode"].(string); ok && sm != "" {
			switch sm {
			case "FILL", "FIT", "CROP", "TILE":
			default:
				return fmt.Sprintf("scaleMode must be FILL, FIT, CROP, or TILE, got: %s", sm)
			}
		}
		if pid, ok := params["parentId"].(string); ok && pid != "" && !ValidNodeID(pid) {
			return fmt.Sprintf("parentId must use colon format e.g. 4029:12345, got: %s", pid)
		}

	case "set_layout_properties":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if !hasAnyParam(params,
			"layoutMode", "layoutWrap",
			"primaryAxisAlignItems", "counterAxisAlignItems",
			"primaryAxisSizingMode", "counterAxisSizingMode",
			"layoutSizingHorizontal", "layoutSizingVertical",
			"layoutAlign", "layoutGrow", "layoutPositioning",
			"paddingTop", "paddingRight", "paddingBottom", "paddingLeft",
			"itemSpacing", "counterAxisSpacing", "clipsContent",
		) {
			return "at least one layout property is required"
		}
		if v, ok := params["layoutMode"].(string); ok && v != "" {
			switch v {
			case "HORIZONTAL", "VERTICAL", "NONE":
			default:
				return fmt.Sprintf("layoutMode must be HORIZONTAL, VERTICAL, or NONE, got: %s", v)
			}
		}
		if v, ok := params["layoutWrap"].(string); ok && v != "" {
			switch v {
			case "NO_WRAP", "WRAP":
			default:
				return fmt.Sprintf("layoutWrap must be NO_WRAP or WRAP, got: %s", v)
			}
		}
		if v, ok := params["primaryAxisAlignItems"].(string); ok && v != "" {
			switch v {
			case "MIN", "CENTER", "MAX", "SPACE_BETWEEN":
			default:
				return fmt.Sprintf("primaryAxisAlignItems must be MIN, CENTER, MAX, or SPACE_BETWEEN, got: %s", v)
			}
		}
		if v, ok := params["counterAxisAlignItems"].(string); ok && v != "" {
			switch v {
			case "MIN", "CENTER", "MAX", "BASELINE":
			default:
				return fmt.Sprintf("counterAxisAlignItems must be MIN, CENTER, MAX, or BASELINE, got: %s", v)
			}
		}
		if v, ok := params["primaryAxisSizingMode"].(string); ok && v != "" {
			switch v {
			case "FIXED", "AUTO":
			default:
				return fmt.Sprintf("primaryAxisSizingMode must be FIXED or AUTO, got: %s", v)
			}
		}
		if v, ok := params["counterAxisSizingMode"].(string); ok && v != "" {
			switch v {
			case "FIXED", "AUTO":
			default:
				return fmt.Sprintf("counterAxisSizingMode must be FIXED or AUTO, got: %s", v)
			}
		}
		if v, ok := params["layoutSizingHorizontal"].(string); ok && v != "" {
			switch v {
			case "FIXED", "HUG", "FILL":
			default:
				return fmt.Sprintf("layoutSizingHorizontal must be FIXED, HUG, or FILL, got: %s", v)
			}
		}
		if v, ok := params["layoutSizingVertical"].(string); ok && v != "" {
			switch v {
			case "FIXED", "HUG", "FILL":
			default:
				return fmt.Sprintf("layoutSizingVertical must be FIXED, HUG, or FILL, got: %s", v)
			}
		}
		if v, ok := params["layoutAlign"].(string); ok && v != "" {
			switch v {
			case "MIN", "CENTER", "MAX", "STRETCH", "INHERIT":
			default:
				return fmt.Sprintf("layoutAlign must be MIN, CENTER, MAX, STRETCH, or INHERIT, got: %s", v)
			}
		}
		if v, ok := params["layoutPositioning"].(string); ok && v != "" {
			switch v {
			case "AUTO", "ABSOLUTE":
			default:
				return fmt.Sprintf("layoutPositioning must be AUTO or ABSOLUTE, got: %s", v)
			}
		}
		if v, ok := params["layoutGrow"].(float64); ok && v != 0 && v != 1 {
			return "layoutGrow must be 0 or 1"
		}

	case "set_text_style":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if !hasAnyParam(params,
			"fontFamily", "fontStyle", "fontSize",
			"textCase", "textAlignHorizontal", "textAlignVertical", "textAutoResize",
			"width", "height",
			"paragraphSpacing", "lineHeight", "letterSpacing",
		) {
			return "at least one text style property is required"
		}
		if v, ok := params["fontSize"].(float64); ok && v <= 0 {
			return "fontSize must be positive"
		}
		if v, ok := params["paragraphSpacing"].(float64); ok && v < 0 {
			return "paragraphSpacing must be non-negative"
		}
		if v, ok := params["width"].(float64); ok && v <= 0 {
			return "width must be positive"
		}
		if v, ok := params["height"].(float64); ok && v <= 0 {
			return "height must be positive"
		}
		if err := validateTextAutoResize(params["textAutoResize"]); err != "" {
			return err
		}
		if v, ok := params["textCase"].(string); ok && v != "" {
			switch v {
			case "ORIGINAL", "UPPER", "LOWER", "TITLE", "SMALL_CAPS", "SMALL_CAPS_FORCED":
			default:
				return fmt.Sprintf("textCase is invalid: %s", v)
			}
		}
		if v, ok := params["textAlignHorizontal"].(string); ok && v != "" {
			switch v {
			case "LEFT", "CENTER", "RIGHT", "JUSTIFIED":
			default:
				return fmt.Sprintf("textAlignHorizontal must be LEFT, CENTER, RIGHT, or JUSTIFIED, got: %s", v)
			}
		}
		if v, ok := params["textAlignVertical"].(string); ok && v != "" {
			switch v {
			case "TOP", "CENTER", "BOTTOM":
			default:
				return fmt.Sprintf("textAlignVertical must be TOP, CENTER, or BOTTOM, got: %s", v)
			}
		}
		if err := validateLineHeight(params["lineHeight"]); err != "" {
			return err
		}
		if err := validateLetterSpacing(params["letterSpacing"]); err != "" {
			return err
		}

	case "set_effects":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		rawEffects, ok := params["effects"]
		if !ok {
			return "effects is required"
		}
		effects, ok := rawEffects.([]interface{})
		if !ok {
			return "effects must be an array"
		}
		for i, raw := range effects {
			effect, ok := raw.(map[string]interface{})
			if !ok {
				return fmt.Sprintf("effects[%d] must be an object", i)
			}
			t, _ := effect["type"].(string)
			switch t {
			case "DROP_SHADOW", "INNER_SHADOW":
				if color, _ := effect["color"].(string); color == "" {
					return fmt.Sprintf("effects[%d].color is required for %s", i, t)
				}
				if radius, ok := effect["radius"].(float64); ok && radius < 0 {
					return fmt.Sprintf("effects[%d].radius must be non-negative", i)
				}
				if spread, ok := effect["spread"].(float64); ok && spread < 0 {
					return fmt.Sprintf("effects[%d].spread must be non-negative", i)
				}
			case "LAYER_BLUR", "BACKGROUND_BLUR":
				if radius, ok := effect["radius"].(float64); !ok || radius < 0 {
					return fmt.Sprintf("effects[%d].radius is required and must be non-negative", i)
				}
			default:
				return fmt.Sprintf("effects[%d].type is invalid: %s", i, t)
			}
		}

	case "apply_styles":
		if len(nodeIDs) == 0 || nodeIDs[0] == "" {
			return "nodeId is required"
		}
		if !ValidNodeID(nodeIDs[0]) {
			return fmt.Sprintf("nodeId must use colon format e.g. 4029:12345, got: %s", nodeIDs[0])
		}
		if !hasAnyNonEmptyString(params, "fillStyleId", "strokeStyleId", "effectStyleId", "textStyleId") {
			return "at least one style id is required"
		}

	case "create_instance":
		if !hasAnyNonEmptyString(params, "componentId", "componentKey", "name") {
			return "one of componentId, componentKey, or name is required"
		}
		if id, ok := params["componentId"].(string); ok && id != "" && !ValidNodeID(id) {
			return fmt.Sprintf("componentId must use colon format e.g. 4029:12345, got: %s", id)
		}
		if pid, ok := params["parentId"].(string); ok && pid != "" && !ValidNodeID(pid) {
			return fmt.Sprintf("parentId must use colon format e.g. 4029:12345, got: %s", pid)
		}
		if props, ok := params["componentProperties"]; ok {
			if _, ok := props.(map[string]interface{}); !ok {
				return "componentProperties must be an object"
			}
		}
	}

	return ""
}

func validExportFormat(f string) bool {
	switch f {
	case "PNG", "SVG", "JPG", "PDF":
		return true
	}
	return false
}

func hasAnyParam(params map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := params[key]; ok {
			return true
		}
	}
	return false
}

func hasAnyNonEmptyString(params map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if v, ok := params[key].(string); ok && v != "" {
			return true
		}
	}
	return false
}

func validateLineHeight(raw interface{}) string {
	if raw == nil {
		return ""
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return "lineHeight must be an object"
	}
	unit, _ := obj["unit"].(string)
	switch unit {
	case "AUTO":
		return ""
	case "PIXELS", "PERCENT", "FONT_SIZE_%":
		if value, ok := obj["value"].(float64); !ok || value < 0 {
			return "lineHeight.value must be non-negative"
		}
		return ""
	default:
		return fmt.Sprintf("lineHeight.unit is invalid: %s", unit)
	}
}

func validateLetterSpacing(raw interface{}) string {
	if raw == nil {
		return ""
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return "letterSpacing must be an object"
	}
	unit, _ := obj["unit"].(string)
	switch unit {
	case "PIXELS", "PERCENT":
		if _, ok := obj["value"].(float64); !ok {
			return "letterSpacing.value is required"
		}
		return ""
	default:
		return fmt.Sprintf("letterSpacing.unit is invalid: %s", unit)
	}
}

func validateTextAutoResize(raw interface{}) string {
	if raw == nil {
		return ""
	}
	value, ok := raw.(string)
	if !ok {
		return "textAutoResize must be a string"
	}
	switch value {
	case "NONE", "WIDTH_AND_HEIGHT", "HEIGHT", "TRUNCATE":
		return ""
	default:
		return fmt.Sprintf("textAutoResize must be NONE, WIDTH_AND_HEIGHT, HEIGHT, or TRUNCATE, got: %s", value)
	}
}
