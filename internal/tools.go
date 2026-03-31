package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all MCP tools on the server.
func RegisterTools(s *server.MCPServer, node *Node) {
	registerReadTools(s, node)
	registerWriteTools(s, node)
}

// RegisterPrompts registers MCP prompts on the server.
func RegisterPrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("read_design_strategy",
		mcp.WithPromptDescription("Best practices for reading Figma designs with figma-mcp-go"),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"Best practices for reading Figma designs",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(`To effectively read a Figma design with figma-mcp-go:

1. Start with get_metadata — understand file name, pages, and current page
2. Use get_pages to list all pages without loading their full trees
3. Use get_design_context (depth=2, detail=compact) for a token-efficient summary of the current selection or page
   - detail=minimal: id/name/type/bounds only (~5% tokens)
   - detail=compact: + fills/strokes/opacity (~30% tokens)
   - detail=full: everything, default (100% tokens)
4. Use search_nodes to find nodes by name or type without dumping the entire tree
5. Drill into specific nodes with get_node or get_nodes_info (prefer batch over single calls)
6. For text-heavy components, use scan_text_nodes to collect all copy at once
7. Use scan_nodes_by_types to find all FRAME/COMPONENT/INSTANCE nodes in a subtree
8. Call get_styles and get_variable_defs once per session to understand the design system
9. Call get_fonts to understand typography usage across the page at a glance
10. Use get_viewport to see what the user is currently looking at in the canvas
11. Use get_reactions to inspect prototype interactions on a node
12. Call get_screenshot last and only when visual confirmation is needed — it is expensive
13. Node IDs use colon format: 4029:12345 — never use hyphens
14. get_local_components now includes componentSets and variantProperties for variant-aware inspection`),
				),
			},
		), nil
	})

	s.AddPrompt(mcp.NewPrompt("design_strategy",
		mcp.WithPromptDescription("Best practices for working with Figma designs"),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"Best practices for working with Figma designs",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(`When working with Figma designs, follow these best practices:

1. Start with Document Structure:
   - First use get_metadata() to understand the current document
   - Use get_pages() to list all pages
   - Plan your layout hierarchy before creating elements
   - Create a main container frame for each screen/section

2. Naming Conventions:
   - Use descriptive, semantic names for all elements
   - Follow a consistent naming pattern (e.g., "Login Screen", "Logo Container", "Email Input")
   - Group related elements with meaningful names

3. Layout Hierarchy:
   - Create parent frames first, then add child elements
   - For forms/login screens:
     * Start with the main screen container frame
     * Create a logo container at the top
     * Group input fields in their own containers
     * Place action buttons (login, submit) after inputs
     * Add secondary elements (forgot password, signup links) last

4. Input Fields Structure:
   - Create a container frame for each input field
   - Include a label text above or inside the input
   - Group related inputs (e.g., username/password) together

5. Element Creation:
   - Use create_frame() for containers and input fields
   - Use create_text() for labels, buttons text, and links
   - When creating multi-line or wrapped text, set width and textAutoResize intentionally
   - Set appropriate colors and styles:
     * Use fillColor for backgrounds
     * Use set_strokes() for borders
     * Set proper fontStyle for different text elements
     * Use set_effects() only when subtle depth or blur improves hierarchy

6. Modifying existing elements:
   - Use set_text() to modify text content of a TEXT node
   - Use set_fills() to change background/fill colors
   - Use move_nodes() / resize_nodes() for position and size adjustments on non-TEXT nodes
   - Use set_text_style() for TEXT sizing, line height, letter spacing, alignment, and textAutoResize

7. Visual Hierarchy:
   - Position elements in logical reading order (top to bottom)
   - Maintain consistent spacing between elements
   - Use appropriate font sizes for different text types:
     * Larger for headings/welcome text
     * Medium for input labels
     * Standard for button text
     * Smaller for helper text/links

8. Best Practices:
   - Verify each creation with get_node()
   - Use parentId to maintain proper hierarchy
   - Group related elements together in frames
   - Keep consistent spacing and alignment
   - All write operations are undoable via Ctrl/Cmd+Z in Figma

Example Login Screen Structure:
- Login Screen (main frame)
  - Logo Container (frame)
    - Logo (text)
  - Welcome Text (text)
  - Input Container (frame)
    - Email Input (frame)
      - Email Label (text)
      - Email Field (frame)
    - Password Input (frame)
      - Password Label (text)
      - Password Field (frame)
  - Login Button (frame)
    - Button Text (text)
  - Helper Links (frame)
    - Forgot Password (text)
    - Don't have account (text)`),
				),
			},
		), nil
	})

	s.AddPrompt(mcp.NewPrompt("text_replacement_strategy",
		mcp.WithPromptDescription("Systematic approach for replacing text in Figma designs"),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"Systematic approach for replacing text in Figma designs",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(`# Intelligent Text Replacement Strategy

## 1. Analyze Design & Identify Structure
- Scan text nodes to understand the overall structure of the design
- Use AI pattern recognition to identify logical groupings:
  * Tables (rows, columns, headers, cells)
  * Lists (items, headers, nested lists)
  * Card groups (similar cards with recurring text fields)
  * Forms (labels, input fields, validation text)
  * Navigation (menu items, breadcrumbs)

scan_text_nodes(nodeId: "node-id")
get_node(nodeId: "node-id")  // optional for extra context

## 2. Strategic Chunking for Complex Designs
- Divide replacement tasks into logical content chunks based on design structure
- Use one of these chunking strategies that best fits the design:
  * Structural Chunking: Table rows/columns, list sections, card groups
  * Spatial Chunking: Top-to-bottom, left-to-right in screen areas
  * Semantic Chunking: Content related to the same topic or functionality
  * Component-Based Chunking: Process similar component instances together

## 3. Progressive Replacement with Verification
- Create a safe copy of the node before bulk replacements
- Replace text chunk by chunk with continuous progress updates
- After each chunk is processed:
  * Export that section with get_screenshot for visual verification
  * Verify text fits properly and maintains design integrity
  * Fix issues before proceeding to the next chunk

// Clone the node to create a safe copy
clone_node(nodeId: "selected-node-id", x: newX, y: newY)

// Replace text one node at a time or in batches
set_text(nodeId: "node-id", text: "New text")

// Verify chunk with targeted image export
get_screenshot(nodeIds: ["chunk-node-id"], format: "PNG", scale: 0.5)

## 4. Intelligent Handling for Table Data
- For tabular content:
  * Process one row or column at a time
  * Maintain alignment and spacing between cells
  * Consider conditional formatting based on cell content
  * Preserve header/data relationships

## 5. Smart Text Adaptation
- Adaptively handle text based on container constraints:
  * Auto-detect space constraints and adjust text length
  * Apply line breaks at appropriate linguistic points
  * Maintain text hierarchy and emphasis

## 6. Final Verification & Context-Aware QA
- After all chunks are processed:
  * Export the entire design at reduced scale for final verification
  * Check for cross-chunk consistency issues
  * Verify proper text flow between different sections
  * Ensure design harmony across the full composition

## 7. Chunk-Specific Export Scale Guidelines
- Scale exports appropriately based on chunk size:
  * Small chunks (1-5 elements): scale 1.0
  * Medium chunks (6-20 elements): scale 0.7
  * Large chunks (21-50 elements): scale 0.5
  * Very large chunks (50+ elements): scale 0.3
  * Full design verification: scale 0.2

## Best Practices
- Preserve Design Intent: Always prioritize design integrity
- Structural Consistency: Maintain alignment, spacing, and hierarchy
- Visual Feedback: Verify each chunk visually before proceeding
- Incremental Improvement: Learn from each chunk to improve subsequent ones
- Respect Content Relationships: Keep related content consistent across chunks`),
				),
			},
		), nil
	})

	s.AddPrompt(mcp.NewPrompt("annotation_conversion_strategy",
		mcp.WithPromptDescription("Strategy for converting manual annotations to Figma's native annotations"),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"Strategy for converting manual annotations to Figma's native annotations",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(`# Automatic Annotation Conversion

## Process Overview
Convert manual annotations (numbered/alphabetical indicators with connected descriptions) to Figma's native annotations:

1. Get selected frame/component information
2. Scan and collect all annotation text nodes
3. Scan target UI elements (components, instances, frames)
4. Match annotations to appropriate UI elements
5. Apply native Figma annotations

## Step 1: Get Selection and Initial Setup

// Get the selected frame/component
get_selection()
// Note the selected node ID, then:
get_annotations(nodeId: "selected-node-id")

## Step 2: Scan Annotation Text Nodes

// Get all text nodes in the selection
scan_text_nodes(nodeId: "selected-node-id")

// Filter and group annotation markers and descriptions
// Markers typically have these characteristics:
// - Short text content (usually single digit/letter)
// - Specific font styles (often bold)
// - Located in a container with "Marker" or "Dot" in the name
// - Have a clear naming pattern (e.g., "1", "2", "3" or "A", "B", "C")

## Step 3: Scan Target UI Elements

// Get all potential target elements that annotations might refer to
scan_nodes_by_types(nodeId: "selected-node-id", types: ["COMPONENT", "INSTANCE", "FRAME"])

## Step 4: Match Annotations to Targets

Match each annotation to its target UI element using these strategies in order of priority:

1. Path-Based Matching:
   - Look at the marker's parent container name in the Figma layer hierarchy
   - Remove any "Marker:" or "Annotation:" prefixes from the parent name
   - Find UI elements that share the same parent name or have it in their path

2. Name-Based Matching:
   - Extract key terms from the annotation description
   - Look for UI elements whose names contain these key terms
   - Particularly effective for form fields, buttons, and labeled components

3. Proximity-Based Matching (fallback):
   - Calculate the center point of the marker using its bounds
   - Find the closest UI element by measuring distances to element centers
   - Use this method when other matching strategies fail

## Step 5: Verify Results

After converting annotations, verify with:
get_annotations(nodeId: "selected-node-id")
get_screenshot(nodeIds: ["selected-node-id"], format: "PNG", scale: 0.5)

This strategy focuses on practical implementation based on real-world usage patterns,
emphasizing the importance of handling various UI elements as annotation targets.`),
				),
			},
		), nil
	})

	s.AddPrompt(mcp.NewPrompt("swap_overrides_instances",
		mcp.WithPromptDescription("Strategy for transferring overrides between component instances in Figma"),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"Strategy for transferring overrides between component instances in Figma",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(`# Swap Component Instance and Override Strategy

## Overview
Transfer content and property overrides from a source instance to one or more target instances
in Figma, maintaining design consistency while reducing manual work.

## Step-by-Step Process

### 1. Selection Analysis
- Use get_selection() to identify the parent component or selected instances
- For parent components, scan for instances with:
  scan_nodes_by_types(nodeId: "parent-id", types: ["INSTANCE"])
- Identify custom slots by name patterns (e.g. "Custom Slot*" or "Instance Slot")
- Determine which is the source instance (with content to copy) and which are targets

### 2. Inspect Source Instance
- Use get_node(nodeId: "source-instance-id") to examine the source instance structure
- Use get_nodes_info(nodeIds: [...]) to batch-inspect multiple instances
- Use scan_text_nodes(nodeId: "source-instance-id") to capture all text content

### 3. Apply Overrides to Targets
- For text overrides: use set_text(nodeId: "target-text-node-id", text: "copied text")
- For fill overrides: use set_fills(nodeId: "target-node-id", color: "#hexcolor")
- For stroke overrides: use set_strokes(nodeId: "target-node-id", color: "#hexcolor")
- Process targets one at a time or identify patterns to apply systematically

### 4. Verification
- Verify results with get_node() or get_design_context()
- Confirm text content and style overrides have transferred successfully
- Use get_screenshot() for visual confirmation if needed

## Key Tips
- Use scan_nodes_by_types to enumerate all instances before starting
- When working with multiple targets, check the full selection with get_selection()
- Prefer reading the full node tree of the source first to understand its structure
- Keep related content consistent across all target instances`),
				),
			},
		), nil
	})

	s.AddPrompt(mcp.NewPrompt("reaction_to_connector_strategy",
		mcp.WithPromptDescription("Strategy for analyzing Figma prototype reactions and mapping interaction flows"),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"Strategy for analyzing Figma prototype reactions and mapping interaction flows",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(`# Strategy: Analyze Figma Prototype Reactions and Map Interaction Flows

## Goal
Process the JSON output from the get_reactions tool to understand prototype flows
and produce a clear, structured map of interactions between screens/nodes.

## Input Data
You will receive JSON data from get_reactions. Each node may contain reactions like:
{
  "trigger": { "type": "ON_CLICK" },
  "action": {
    "type": "NAVIGATE",
    "destinationId": "destination-node-id"
  }
}

## Step-by-Step Process

### 1. Gather Context
- Call get_nodes_info(nodeIds: [...]) on all relevant nodes to get their names and types
- Call get_design_context(depth: 2, detail: "minimal") to understand the page structure

### 2. Filter and Transform Reactions
- Iterate through the get_reactions JSON output
- Keep only reactions where action type implies navigation:
  * NAVIGATE, OPEN_OVERLAY, SWAP_OVERLAY
  * Ignore: CHANGE_TO, CLOSE_OVERLAY, and others without a destinationId
- Extract per reaction:
  * sourceNodeId: the node the reaction belongs to
  * destinationId: action.destinationId
  * actionType: action.type
  * triggerType: trigger.type

### 3. Generate Flow Map
For each valid reaction, create a human-readable description:
- "On click → navigate to [Destination Name]"
- "On drag → open [Destination Name] overlay"
- "On hover → swap to [Destination Name]"

Combine these into a structured flow map grouped by source screen.

### 4. Output Format
Produce a summary like:

Flow Map:
- [Screen A] --ON_CLICK/NAVIGATE--> [Screen B]
- [Screen A] --ON_CLICK/OPEN_OVERLAY--> [Modal C]
- [Screen B] --ON_CLICK/NAVIGATE--> [Screen C]

### 5. Verification
- Use get_screenshot(nodeIds: [...]) on key screens to visually confirm the flow
- Cross-check node names from get_nodes_info with the flow map

## Notes
- Node IDs use colon format: 4029:12345 — never use hyphens
- Use get_reactions on a set of nodes that represent screens or interactive frames
- Focus on NAVIGATE actions for the primary user journey`),
				),
			},
		), nil
	})
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// makeHandler creates a simple tool handler with no parameters.
func makeHandler(node *Node, command string, nodeIDs []string, params map[string]interface{}) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resp, err := node.Send(ctx, command, nodeIDs, params)
		return renderResponse(resp, err)
	}
}

// renderResponse converts a BridgeResponse into an MCP tool result.
func renderResponse(resp BridgeResponse, err error) (*mcp.CallToolResult, error) {
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if resp.Error != "" {
		return mcp.NewToolResultError(resp.Error), nil
	}
	text, err := json.Marshal(resp.Data)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal response: %v", err)), nil
	}
	return mcp.NewToolResultText(string(text)), nil
}

// toStringSlice converts []interface{} to []string.
func toStringSlice(raw []interface{}) []string {
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// ── save_screenshots ─────────────────────────────────────────────────────────

type saveItem struct {
	NodeID     string  `json:"nodeId"`
	OutputPath string  `json:"outputPath"`
	Format     string  `json:"format,omitempty"`
	Scale      float64 `json:"scale,omitempty"`
}

type saveResult struct {
	Index        int     `json:"index"`
	NodeID       string  `json:"nodeId"`
	NodeName     string  `json:"nodeName,omitempty"`
	OutputPath   string  `json:"outputPath"`
	Format       string  `json:"format,omitempty"`
	Width        float64 `json:"width,omitempty"`
	Height       float64 `json:"height,omitempty"`
	BytesWritten int     `json:"bytesWritten,omitempty"`
	Success      bool    `json:"success"`
	Error        string  `json:"error,omitempty"`
}

func executeSaveScreenshots(ctx context.Context, node *Node, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rawItems, _ := req.GetArguments()["items"].([]interface{})
	defaultFormat, _ := req.GetArguments()["format"].(string)
	defaultScale, _ := req.GetArguments()["scale"].(float64)

	workDir, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("getwd: %v", err)), nil
	}

	results := make([]saveResult, 0, len(rawItems))
	succeeded, failed := 0, 0

	for i, rawItem := range rawItems {
		item, err := parseSaveItem(rawItem)
		if err != nil {
			results = append(results, saveResult{Index: i, Error: err.Error()})
			failed++
			continue
		}

		r := saveScreenshotItem(ctx, node, item, i, workDir, defaultFormat, defaultScale)
		results = append(results, r)
		if r.Success {
			succeeded++
		} else {
			failed++
		}
	}

	out, err := json.Marshal(map[string]interface{}{
		"total":     len(results),
		"succeeded": succeeded,
		"failed":    failed,
		"hasErrors": failed > 0,
		"results":   results,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal results: %v", err)), nil
	}
	return mcp.NewToolResultText(string(out)), nil
}

func saveScreenshotItem(ctx context.Context, node *Node, item saveItem, index int, workDir, defaultFormat string, defaultScale float64) saveResult {
	resolvedPath, err := resolveOutputPath(item.OutputPath, workDir)
	if err != nil {
		return saveResult{Index: index, NodeID: item.NodeID, OutputPath: item.OutputPath, Error: err.Error()}
	}

	format := coalesce(item.Format, defaultFormat)
	inferredFormat := inferFormat(resolvedPath)
	if format == "" {
		format = inferredFormat
	}
	if format == "" {
		format = "PNG"
	}
	if inferredFormat != "" && format != inferredFormat {
		return saveResult{Index: index, NodeID: item.NodeID, OutputPath: resolvedPath,
			Error: fmt.Sprintf("format %s conflicts with file extension %s", format, inferredFormat)}
	}

	scale := item.Scale
	if scale <= 0 {
		scale = defaultScale
	}

	params := map[string]interface{}{"format": format}
	if scale > 0 {
		params["scale"] = scale
	}

	resp, err := node.Send(ctx, "get_screenshot", []string{item.NodeID}, params)
	if err != nil {
		return saveResult{Index: index, NodeID: item.NodeID, OutputPath: resolvedPath, Error: err.Error()}
	}
	if resp.Error != "" {
		return saveResult{Index: index, NodeID: item.NodeID, OutputPath: resolvedPath, Error: resp.Error}
	}

	export, err := extractScreenshotExport(resp.Data)
	if err != nil {
		return saveResult{Index: index, NodeID: item.NodeID, OutputPath: resolvedPath, Error: err.Error()}
	}

	bytes, err := writeBase64(export.Base64, resolvedPath)
	if err != nil {
		return saveResult{Index: index, NodeID: item.NodeID, OutputPath: resolvedPath, Error: err.Error()}
	}

	return saveResult{
		Index:        index,
		NodeID:       export.NodeID,
		NodeName:     export.NodeName,
		OutputPath:   resolvedPath,
		Format:       format,
		Width:        export.Width,
		Height:       export.Height,
		BytesWritten: bytes,
		Success:      true,
	}
}

type screenshotExport struct {
	NodeID   string  `json:"nodeId"`
	NodeName string  `json:"nodeName"`
	Base64   string  `json:"base64"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
}

func extractScreenshotExport(data interface{}) (screenshotExport, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return screenshotExport{}, err
	}
	var wrapper struct {
		Exports []screenshotExport `json:"exports"`
	}
	if err := json.Unmarshal(b, &wrapper); err != nil {
		return screenshotExport{}, err
	}
	if len(wrapper.Exports) == 0 {
		return screenshotExport{}, errors.New("no screenshot export returned by plugin")
	}
	return wrapper.Exports[0], nil
}

func writeBase64(b64, outputPath string) (int, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return 0, fmt.Errorf("base64 decode: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return 0, fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return 0, fmt.Errorf("file already exists at outputPath: %s", outputPath)
		}
		return 0, err
	}
	defer f.Close()
	n, err := f.Write(data)
	return n, err
}

func resolveOutputPath(outputPath, workDir string) (string, error) {
	// Reject absolute paths early — filepath.Join on Windows can let a drive-letter
	// absolute path escape the workDir check via filepath.Rel.
	if filepath.IsAbs(outputPath) {
		return "", fmt.Errorf("outputPath must be a relative path, got: %s", outputPath)
	}

	resolved := filepath.Join(workDir, outputPath)
	rel, err := filepath.Rel(workDir, resolved)
	if err != nil {
		return "", fmt.Errorf("outputPath must be inside the working directory: %s", workDir)
	}

	// Convert to forward slashes before prefix check so Windows paths like
	// "C:\.." don't bypass the ".." detection.
	if strings.HasPrefix(filepath.ToSlash(rel), "..") {
		return "", fmt.Errorf("outputPath must be inside the working directory: %s", workDir)
	}

	return resolved, nil
}

func inferFormat(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return "PNG"
	case ".svg":
		return "SVG"
	case ".jpg", ".jpeg":
		return "JPG"
	case ".pdf":
		return "PDF"
	}
	return ""
}

func parseSaveItem(raw interface{}) (saveItem, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return saveItem{}, err
	}
	var item saveItem
	if err := json.Unmarshal(b, &item); err != nil {
		return saveItem{}, err
	}
	return item, nil
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
