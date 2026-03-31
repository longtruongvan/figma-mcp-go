package internal

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerWriteTools(s *server.MCPServer, node *Node) {
	// ── Write — Create ───────────────────────────────────────────────────

	s.AddTool(mcp.NewTool("create_frame",
		mcp.WithDescription("Create a new frame on the current page or inside a parent node."),
		mcp.WithNumber("x", mcp.Description("X position (default 0)")),
		mcp.WithNumber("y", mcp.Description("Y position (default 0)")),
		mcp.WithNumber("width", mcp.Description("Width in pixels (default 100)")),
		mcp.WithNumber("height", mcp.Description("Height in pixels (default 100)")),
		mcp.WithString("name", mcp.Description("Frame name")),
		mcp.WithString("fillColor", mcp.Description("Fill color as hex e.g. #FFFFFF")),
		mcp.WithString("layoutMode", mcp.Description("Auto-layout direction: HORIZONTAL, VERTICAL, or NONE")),
		mcp.WithNumber("paddingTop", mcp.Description("Auto-layout top padding")),
		mcp.WithNumber("paddingRight", mcp.Description("Auto-layout right padding")),
		mcp.WithNumber("paddingBottom", mcp.Description("Auto-layout bottom padding")),
		mcp.WithNumber("paddingLeft", mcp.Description("Auto-layout left padding")),
		mcp.WithNumber("itemSpacing", mcp.Description("Auto-layout gap between children")),
		mcp.WithString("parentId", mcp.Description("Parent node ID in colon format. Defaults to current page.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := req.GetArguments()
		resp, err := node.Send(ctx, "create_frame", nil, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("create_rectangle",
		mcp.WithDescription("Create a new rectangle on the current page or inside a parent node."),
		mcp.WithNumber("x", mcp.Description("X position (default 0)")),
		mcp.WithNumber("y", mcp.Description("Y position (default 0)")),
		mcp.WithNumber("width", mcp.Description("Width in pixels (default 100)")),
		mcp.WithNumber("height", mcp.Description("Height in pixels (default 100)")),
		mcp.WithString("name", mcp.Description("Rectangle name")),
		mcp.WithString("fillColor", mcp.Description("Fill color as hex e.g. #FF5733")),
		mcp.WithNumber("cornerRadius", mcp.Description("Corner radius in pixels")),
		mcp.WithString("parentId", mcp.Description("Parent node ID in colon format. Defaults to current page.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := req.GetArguments()
		resp, err := node.Send(ctx, "create_rectangle", nil, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("create_ellipse",
		mcp.WithDescription("Create a new ellipse (circle/oval) on the current page or inside a parent node."),
		mcp.WithNumber("x", mcp.Description("X position (default 0)")),
		mcp.WithNumber("y", mcp.Description("Y position (default 0)")),
		mcp.WithNumber("width", mcp.Description("Width in pixels (default 100)")),
		mcp.WithNumber("height", mcp.Description("Height in pixels (default 100)")),
		mcp.WithString("name", mcp.Description("Ellipse name")),
		mcp.WithString("fillColor", mcp.Description("Fill color as hex e.g. #3B82F6")),
		mcp.WithString("parentId", mcp.Description("Parent node ID in colon format. Defaults to current page.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := req.GetArguments()
		resp, err := node.Send(ctx, "create_ellipse", nil, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("create_text",
		mcp.WithDescription("Create a new text node on the current page or inside a parent node, with optional typography and text-box sizing controls."),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text content"),
		),
		mcp.WithNumber("x", mcp.Description("X position (default 0)")),
		mcp.WithNumber("y", mcp.Description("Y position (default 0)")),
		mcp.WithNumber("fontSize", mcp.Description("Font size in pixels (default 14)")),
		mcp.WithString("fontFamily", mcp.Description("Font family e.g. Inter (default Inter)")),
		mcp.WithString("fontStyle", mcp.Description("Font style e.g. Regular, Bold (default Regular)")),
		mcp.WithString("textAutoResize",
			mcp.Description("Text box sizing mode: NONE, WIDTH_AND_HEIGHT, HEIGHT, or TRUNCATE"),
			mcp.Enum("NONE", "WIDTH_AND_HEIGHT", "HEIGHT", "TRUNCATE"),
		),
		mcp.WithNumber("width", mcp.Description("Text box width in pixels. If provided alone, the node defaults to HEIGHT auto-resize for wrapping.")),
		mcp.WithNumber("height", mcp.Description("Text box height in pixels. If provided, the node defaults to NONE auto-resize unless textAutoResize is explicitly set.")),
		mcp.WithString("textAlignHorizontal",
			mcp.Description("Horizontal text alignment"),
			mcp.Enum("LEFT", "CENTER", "RIGHT", "JUSTIFIED"),
		),
		mcp.WithString("textAlignVertical",
			mcp.Description("Vertical text alignment"),
			mcp.Enum("TOP", "CENTER", "BOTTOM"),
		),
		mcp.WithNumber("paragraphSpacing", mcp.Description("Paragraph spacing")),
		mcp.WithObject("lineHeight",
			mcp.Description("Line height object: {unit: AUTO} or {unit: PIXELS|PERCENT|FONT_SIZE_%, value: number}"),
			mcp.Properties(map[string]interface{}{
				"unit": map[string]interface{}{"type": "string"},
				"value": map[string]interface{}{"type": "number"},
			}),
		),
		mcp.WithObject("letterSpacing",
			mcp.Description("Letter spacing object: {unit: PIXELS|PERCENT, value: number}"),
			mcp.Properties(map[string]interface{}{
				"unit": map[string]interface{}{"type": "string"},
				"value": map[string]interface{}{"type": "number"},
			}),
		),
		mcp.WithString("fillColor", mcp.Description("Text color as hex e.g. #000000")),
		mcp.WithString("name", mcp.Description("Node name")),
		mcp.WithString("parentId", mcp.Description("Parent node ID in colon format. Defaults to current page.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := map[string]interface{}{
			"text": req.GetArguments()["text"],
		}
		for _, key := range []string{
			"fontFamily", "fontStyle", "fillColor", "name", "parentId",
			"textAutoResize", "textAlignHorizontal", "textAlignVertical",
		} {
			if v, ok := req.GetArguments()[key].(string); ok && v != "" {
				params[key] = v
			}
		}
		for _, key := range []string{"x", "y", "fontSize", "width", "height", "paragraphSpacing"} {
			if v, ok := req.GetArguments()[key].(float64); ok {
				params[key] = v
			}
		}
		if v, ok := req.GetArguments()["lineHeight"].(map[string]interface{}); ok {
			params["lineHeight"] = v
		}
		if v, ok := req.GetArguments()["letterSpacing"].(map[string]interface{}); ok {
			params["letterSpacing"] = v
		}
		resp, err := node.Send(ctx, "create_text", nil, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("import_image",
		mcp.WithDescription("Import a base64-encoded image into Figma as a rectangle with an image fill. Use get_screenshot to capture images or provide your own base64 PNG/JPG."),
		mcp.WithString("imageData",
			mcp.Required(),
			mcp.Description("Base64-encoded image data (PNG or JPG)"),
		),
		mcp.WithNumber("x", mcp.Description("X position (default 0)")),
		mcp.WithNumber("y", mcp.Description("Y position (default 0)")),
		mcp.WithNumber("width", mcp.Description("Width in pixels (default 200)")),
		mcp.WithNumber("height", mcp.Description("Height in pixels (default 200)")),
		mcp.WithString("name", mcp.Description("Node name")),
		mcp.WithString("scaleMode", mcp.Description("Image scale mode: FILL (default), FIT, CROP, or TILE")),
		mcp.WithString("parentId", mcp.Description("Parent node ID in colon format. Defaults to current page.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := map[string]interface{}{
			"imageData": req.GetArguments()["imageData"],
		}
		if x, ok := req.GetArguments()["x"].(float64); ok {
			params["x"] = x
		}
		if y, ok := req.GetArguments()["y"].(float64); ok {
			params["y"] = y
		}
		if w, ok := req.GetArguments()["width"].(float64); ok {
			params["width"] = w
		}
		if h, ok := req.GetArguments()["height"].(float64); ok {
			params["height"] = h
		}
		if n, ok := req.GetArguments()["name"].(string); ok && n != "" {
			params["name"] = n
		}
		if sm, ok := req.GetArguments()["scaleMode"].(string); ok && sm != "" {
			params["scaleMode"] = sm
		}
		if pid, ok := req.GetArguments()["parentId"].(string); ok && pid != "" {
			params["parentId"] = pid
		}
		resp, err := node.Send(ctx, "import_image", nil, params)
		return renderResponse(resp, err)
	})

	// ── Write — Modify ───────────────────────────────────────────────────

	s.AddTool(mcp.NewTool("set_text",
		mcp.WithDescription("Update the text content of an existing TEXT node."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("TEXT node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("New text content"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		text, _ := req.GetArguments()["text"].(string)
		resp, err := node.Send(ctx, "set_text", []string{nodeID}, map[string]interface{}{"text": text})
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("set_fills",
		mcp.WithDescription("Set the fill color of a node."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("color",
			mcp.Required(),
			mcp.Description("Fill color as hex e.g. #FF5733"),
		),
		mcp.WithNumber("opacity", mcp.Description("Fill opacity 0–1 (default 1)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{
			"color": req.GetArguments()["color"],
		}
		if op, ok := req.GetArguments()["opacity"].(float64); ok {
			params["opacity"] = op
		}
		resp, err := node.Send(ctx, "set_fills", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("set_strokes",
		mcp.WithDescription("Set the stroke color and weight of a node."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("color",
			mcp.Required(),
			mcp.Description("Stroke color as hex e.g. #000000"),
		),
		mcp.WithNumber("strokeWeight", mcp.Description("Stroke weight in pixels (default 1)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{
			"color": req.GetArguments()["color"],
		}
		if sw, ok := req.GetArguments()["strokeWeight"].(float64); ok {
			params["strokeWeight"] = sw
		}
		resp, err := node.Send(ctx, "set_strokes", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("set_layout_properties",
		mcp.WithDescription("Update auto-layout, sizing, spacing, and layout-child properties on a node."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("layoutMode",
			mcp.Description("Auto-layout mode: HORIZONTAL, VERTICAL, or NONE"),
			mcp.Enum("HORIZONTAL", "VERTICAL", "NONE"),
		),
		mcp.WithString("layoutWrap",
			mcp.Description("Wrapping mode for auto-layout frames: NO_WRAP or WRAP"),
			mcp.Enum("NO_WRAP", "WRAP"),
		),
		mcp.WithString("primaryAxisAlignItems",
			mcp.Description("Main-axis alignment: MIN, CENTER, MAX, or SPACE_BETWEEN"),
			mcp.Enum("MIN", "CENTER", "MAX", "SPACE_BETWEEN"),
		),
		mcp.WithString("counterAxisAlignItems",
			mcp.Description("Cross-axis alignment: MIN, CENTER, MAX, or BASELINE"),
			mcp.Enum("MIN", "CENTER", "MAX", "BASELINE"),
		),
		mcp.WithString("primaryAxisSizingMode",
			mcp.Description("Primary axis sizing mode: FIXED or AUTO"),
			mcp.Enum("FIXED", "AUTO"),
		),
		mcp.WithString("counterAxisSizingMode",
			mcp.Description("Counter axis sizing mode: FIXED or AUTO"),
			mcp.Enum("FIXED", "AUTO"),
		),
		mcp.WithString("layoutSizingHorizontal",
			mcp.Description("Horizontal sizing in the Figma UI: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL"),
		),
		mcp.WithString("layoutSizingVertical",
			mcp.Description("Vertical sizing in the Figma UI: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL"),
		),
		mcp.WithString("layoutAlign",
			mcp.Description("Auto-layout child alignment: MIN, CENTER, MAX, STRETCH, or INHERIT"),
			mcp.Enum("MIN", "CENTER", "MAX", "STRETCH", "INHERIT"),
		),
		mcp.WithNumber("layoutGrow",
			mcp.Description("Auto-layout child grow factor: 0 or 1"),
		),
		mcp.WithString("layoutPositioning",
			mcp.Description("Auto-layout child positioning: AUTO or ABSOLUTE"),
			mcp.Enum("AUTO", "ABSOLUTE"),
		),
		mcp.WithNumber("paddingTop", mcp.Description("Top padding")),
		mcp.WithNumber("paddingRight", mcp.Description("Right padding")),
		mcp.WithNumber("paddingBottom", mcp.Description("Bottom padding")),
		mcp.WithNumber("paddingLeft", mcp.Description("Left padding")),
		mcp.WithNumber("itemSpacing", mcp.Description("Spacing between items")),
		mcp.WithNumber("counterAxisSpacing", mcp.Description("Spacing between wrapped rows/columns")),
		mcp.WithBoolean("clipsContent", mcp.Description("Whether overflowing children are clipped")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{}
		for _, key := range []string{
			"layoutMode", "layoutWrap",
			"primaryAxisAlignItems", "counterAxisAlignItems",
			"primaryAxisSizingMode", "counterAxisSizingMode",
			"layoutSizingHorizontal", "layoutSizingVertical",
			"layoutAlign", "layoutPositioning",
		} {
			if v, ok := req.GetArguments()[key].(string); ok && v != "" {
				params[key] = v
			}
		}
		for _, key := range []string{
			"layoutGrow",
			"paddingTop", "paddingRight", "paddingBottom", "paddingLeft",
			"itemSpacing", "counterAxisSpacing",
		} {
			if v, ok := req.GetArguments()[key].(float64); ok {
				params[key] = v
			}
		}
		if v, ok := req.GetArguments()["clipsContent"].(bool); ok {
			params["clipsContent"] = v
		}
		resp, err := node.Send(ctx, "set_layout_properties", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("set_text_style",
		mcp.WithDescription("Update typography and text-box properties on a TEXT node."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("TEXT node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("fontFamily", mcp.Description("Font family e.g. Inter")),
		mcp.WithString("fontStyle", mcp.Description("Font style e.g. Regular, Bold")),
		mcp.WithNumber("fontSize", mcp.Description("Font size in pixels")),
		mcp.WithString("textCase",
			mcp.Description("Text case override"),
			mcp.Enum("ORIGINAL", "UPPER", "LOWER", "TITLE", "SMALL_CAPS", "SMALL_CAPS_FORCED"),
		),
		mcp.WithString("textAlignHorizontal",
			mcp.Description("Horizontal text alignment"),
			mcp.Enum("LEFT", "CENTER", "RIGHT", "JUSTIFIED"),
		),
		mcp.WithString("textAlignVertical",
			mcp.Description("Vertical text alignment"),
			mcp.Enum("TOP", "CENTER", "BOTTOM"),
		),
		mcp.WithString("textAutoResize",
			mcp.Description("Text box sizing mode: NONE, WIDTH_AND_HEIGHT, HEIGHT, or TRUNCATE"),
			mcp.Enum("NONE", "WIDTH_AND_HEIGHT", "HEIGHT", "TRUNCATE"),
		),
		mcp.WithNumber("width", mcp.Description("Text box width in pixels")),
		mcp.WithNumber("height", mcp.Description("Text box height in pixels")),
		mcp.WithNumber("paragraphSpacing", mcp.Description("Paragraph spacing")),
		mcp.WithObject("lineHeight",
			mcp.Description("Line height object: {unit: AUTO} or {unit: PIXELS|PERCENT|FONT_SIZE_%, value: number}"),
			mcp.Properties(map[string]interface{}{
				"unit": map[string]interface{}{"type": "string"},
				"value": map[string]interface{}{"type": "number"},
			}),
		),
		mcp.WithObject("letterSpacing",
			mcp.Description("Letter spacing object: {unit: PIXELS|PERCENT, value: number}"),
			mcp.Properties(map[string]interface{}{
				"unit": map[string]interface{}{"type": "string"},
				"value": map[string]interface{}{"type": "number"},
			}),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{}
		for _, key := range []string{"fontFamily", "fontStyle", "textCase", "textAlignHorizontal", "textAlignVertical", "textAutoResize"} {
			if v, ok := req.GetArguments()[key].(string); ok && v != "" {
				params[key] = v
			}
		}
		for _, key := range []string{"fontSize", "width", "height", "paragraphSpacing"} {
			if v, ok := req.GetArguments()[key].(float64); ok {
				params[key] = v
			}
		}
		if v, ok := req.GetArguments()["lineHeight"].(map[string]interface{}); ok {
			params["lineHeight"] = v
		}
		if v, ok := req.GetArguments()["letterSpacing"].(map[string]interface{}); ok {
			params["letterSpacing"] = v
		}
		resp, err := node.Send(ctx, "set_text_style", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("set_effects",
		mcp.WithDescription("Replace a node's effects with shadows and/or blurs. Pass an empty array to clear all effects."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithArray("effects",
			mcp.Required(),
			mcp.Description("List of effect objects. Supported types: DROP_SHADOW, INNER_SHADOW, LAYER_BLUR, BACKGROUND_BLUR"),
			mcp.Items(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"type": map[string]interface{}{"type": "string"},
					"color": map[string]interface{}{"type": "string"},
					"x": map[string]interface{}{"type": "number"},
					"y": map[string]interface{}{"type": "number"},
					"radius": map[string]interface{}{"type": "number"},
					"spread": map[string]interface{}{"type": "number"},
					"blendMode": map[string]interface{}{"type": "string"},
					"visible": map[string]interface{}{"type": "boolean"},
					"showShadowBehindNode": map[string]interface{}{"type": "boolean"},
				},
			}),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{}
		if raw, ok := req.GetArguments()["effects"].([]interface{}); ok {
			params["effects"] = raw
		}
		resp, err := node.Send(ctx, "set_effects", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("apply_styles",
		mcp.WithDescription("Apply local fill, stroke, effect, and/or text styles to a node using style IDs."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("fillStyleId", mcp.Description("Local paint style ID to apply to fills")),
		mcp.WithString("strokeStyleId", mcp.Description("Local paint style ID to apply to strokes")),
		mcp.WithString("effectStyleId", mcp.Description("Local effect style ID to apply to effects")),
		mcp.WithString("textStyleId", mcp.Description("Local text style ID to apply to a TEXT node")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{}
		for _, key := range []string{"fillStyleId", "strokeStyleId", "effectStyleId", "textStyleId"} {
			if v, ok := req.GetArguments()[key].(string); ok && v != "" {
				params[key] = v
			}
		}
		resp, err := node.Send(ctx, "apply_styles", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("create_instance",
		mcp.WithDescription("Create an instance from a local component using componentId, componentKey, or name."),
		mcp.WithString("componentId", mcp.Description("Local component node ID in colon format")),
		mcp.WithString("componentKey", mcp.Description("Local component key")),
		mcp.WithString("name", mcp.Description("Local component name fallback if componentId/componentKey are unknown")),
		mcp.WithString("parentId", mcp.Description("Parent node ID in colon format. Defaults to current page.")),
		mcp.WithNumber("x", mcp.Description("X position")),
		mcp.WithNumber("y", mcp.Description("Y position")),
		mcp.WithObject("componentProperties",
			mcp.Description("Optional component properties to pass to instance.setProperties()"),
			mcp.AdditionalProperties(map[string]interface{}{
				"anyOf": []map[string]interface{}{
					{"type": "string"},
					{"type": "boolean"},
				},
			}),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := map[string]interface{}{}
		for _, key := range []string{"componentId", "componentKey", "name", "parentId"} {
			if v, ok := req.GetArguments()[key].(string); ok && v != "" {
				params[key] = v
			}
		}
		for _, key := range []string{"x", "y"} {
			if v, ok := req.GetArguments()[key].(float64); ok {
				params[key] = v
			}
		}
		if v, ok := req.GetArguments()["componentProperties"].(map[string]interface{}); ok && len(v) > 0 {
			params["componentProperties"] = v
		}
		resp, err := node.Send(ctx, "create_instance", nil, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("move_nodes",
		mcp.WithDescription("Move one or more nodes to an absolute position on the canvas."),
		mcp.WithArray("nodeIds",
			mcp.Required(),
			mcp.Description("Node IDs in colon format e.g. ['4029:12345']"),
		),
		mcp.WithNumber("x", mcp.Description("Target X position")),
		mcp.WithNumber("y", mcp.Description("Target Y position")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, _ := req.GetArguments()["nodeIds"].([]interface{})
		nodeIDs := toStringSlice(raw)
		params := map[string]interface{}{}
		if x, ok := req.GetArguments()["x"].(float64); ok {
			params["x"] = x
		}
		if y, ok := req.GetArguments()["y"].(float64); ok {
			params["y"] = y
		}
		resp, err := node.Send(ctx, "move_nodes", nodeIDs, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("resize_nodes",
		mcp.WithDescription("Resize one or more nodes. Provide width, height, or both."),
		mcp.WithArray("nodeIds",
			mcp.Required(),
			mcp.Description("Node IDs in colon format e.g. ['4029:12345']"),
		),
		mcp.WithNumber("width", mcp.Description("New width in pixels")),
		mcp.WithNumber("height", mcp.Description("New height in pixels")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, _ := req.GetArguments()["nodeIds"].([]interface{})
		nodeIDs := toStringSlice(raw)
		params := map[string]interface{}{}
		if w, ok := req.GetArguments()["width"].(float64); ok {
			params["width"] = w
		}
		if h, ok := req.GetArguments()["height"].(float64); ok {
			params["height"] = h
		}
		resp, err := node.Send(ctx, "resize_nodes", nodeIDs, params)
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("rename_node",
		mcp.WithDescription("Rename a node."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("New name for the node"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		name, _ := req.GetArguments()["name"].(string)
		resp, err := node.Send(ctx, "rename_node", []string{nodeID}, map[string]interface{}{"name": name})
		return renderResponse(resp, err)
	})

	s.AddTool(mcp.NewTool("clone_node",
		mcp.WithDescription("Clone an existing node, optionally repositioning it or placing it in a new parent."),
		mcp.WithString("nodeId",
			mcp.Required(),
			mcp.Description("Source node ID in colon format e.g. '4029:12345'"),
		),
		mcp.WithNumber("x", mcp.Description("X position of the clone")),
		mcp.WithNumber("y", mcp.Description("Y position of the clone")),
		mcp.WithString("parentId", mcp.Description("Parent node ID for the clone. Defaults to same parent as source.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, _ := req.GetArguments()["nodeId"].(string)
		params := map[string]interface{}{}
		if x, ok := req.GetArguments()["x"].(float64); ok {
			params["x"] = x
		}
		if y, ok := req.GetArguments()["y"].(float64); ok {
			params["y"] = y
		}
		if pid, ok := req.GetArguments()["parentId"].(string); ok && pid != "" {
			params["parentId"] = pid
		}
		resp, err := node.Send(ctx, "clone_node", []string{nodeID}, params)
		return renderResponse(resp, err)
	})

	// ── Write — Delete ───────────────────────────────────────────────────

	s.AddTool(mcp.NewTool("delete_nodes",
		mcp.WithDescription("Delete one or more nodes. This cannot be undone via MCP — use with care."),
		mcp.WithArray("nodeIds",
			mcp.Required(),
			mcp.Description("Node IDs to delete in colon format e.g. ['4029:12345']"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, _ := req.GetArguments()["nodeIds"].([]interface{})
		nodeIDs := toStringSlice(raw)
		resp, err := node.Send(ctx, "delete_nodes", nodeIDs, nil)
		return renderResponse(resp, err)
	})
}
