// Write handlers — all Figma mutation operations.
// Returns null for unknown request types so the caller can surface an error.
// Every mutation calls figma.commitUndo() so changes are undoable via Cmd/Ctrl+Z.

import { getBounds } from "./serializers";
import {
  makeSolidPaint,
  getParentNode,
  base64ToBytes,
  loadAllFonts,
  parseLetterSpacing,
  parseLineHeight,
  makeEffect,
  resolveComponentNode,
  mapComponentProperties,
} from "./write-helpers";

export const handleWriteRequest = async (request: any) => {
  switch (request.type) {
    case "create_frame": {
      const p = request.params || {};
      const parent = await getParentNode(p.parentId);
      const frame = figma.createFrame();
      frame.resize(p.width || 100, p.height || 100);
      frame.x = p.x != null ? p.x : 0;
      frame.y = p.y != null ? p.y : 0;
      if (p.name) frame.name = p.name;
      if (p.fillColor) frame.fills = [makeSolidPaint(p.fillColor)];
      if (p.layoutMode && ["HORIZONTAL", "VERTICAL", "NONE"].includes(p.layoutMode)) {
        frame.layoutMode = p.layoutMode;
      }
      if (p.paddingTop != null) frame.paddingTop = p.paddingTop;
      if (p.paddingRight != null) frame.paddingRight = p.paddingRight;
      if (p.paddingBottom != null) frame.paddingBottom = p.paddingBottom;
      if (p.paddingLeft != null) frame.paddingLeft = p.paddingLeft;
      if (p.itemSpacing != null) frame.itemSpacing = p.itemSpacing;
      (parent as any).appendChild(frame);
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: frame.id, name: frame.name, type: frame.type, bounds: getBounds(frame) },
      };
    }

    case "create_rectangle": {
      const p = request.params || {};
      const parent = await getParentNode(p.parentId);
      const rect = figma.createRectangle();
      rect.resize(p.width || 100, p.height || 100);
      rect.x = p.x != null ? p.x : 0;
      rect.y = p.y != null ? p.y : 0;
      if (p.name) rect.name = p.name;
      if (p.fillColor) rect.fills = [makeSolidPaint(p.fillColor)];
      if (p.cornerRadius != null) rect.cornerRadius = p.cornerRadius;
      (parent as any).appendChild(rect);
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: rect.id, name: rect.name, type: rect.type, bounds: getBounds(rect) },
      };
    }

    case "create_ellipse": {
      const p = request.params || {};
      const parent = await getParentNode(p.parentId);
      const ellipse = figma.createEllipse();
      ellipse.resize(p.width || 100, p.height || 100);
      ellipse.x = p.x != null ? p.x : 0;
      ellipse.y = p.y != null ? p.y : 0;
      if (p.name) ellipse.name = p.name;
      if (p.fillColor) ellipse.fills = [makeSolidPaint(p.fillColor)];
      (parent as any).appendChild(ellipse);
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: ellipse.id, name: ellipse.name, type: ellipse.type, bounds: getBounds(ellipse) },
      };
    }

    case "create_text": {
      const p = request.params || {};
      const parent = await getParentNode(p.parentId);
      const fontFamily = p.fontFamily || "Inter";
      const fontStyle = p.fontStyle || "Regular";
      await figma.loadFontAsync({ family: fontFamily, style: fontStyle });
      const textNode = figma.createText();
      textNode.fontName = { family: fontFamily, style: fontStyle };
      if (p.fontSize) textNode.fontSize = p.fontSize;
      textNode.characters = p.text || "";
      textNode.x = p.x != null ? p.x : 0;
      textNode.y = p.y != null ? p.y : 0;
      if (p.name) textNode.name = p.name;
      if (p.fillColor) textNode.fills = [makeSolidPaint(p.fillColor)];
      (parent as any).appendChild(textNode);
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: textNode.id, name: textNode.name, type: textNode.type, bounds: getBounds(textNode) },
      };
    }

    case "set_text": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId);
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      if (node.type !== "TEXT") throw new Error(`Node ${nodeId} is not a TEXT node`);
      const fontName = typeof node.fontName === "symbol"
        ? { family: "Inter", style: "Regular" }
        : node.fontName;
      await figma.loadFontAsync(fontName);
      node.characters = p.text;
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name, characters: node.characters },
      };
    }

    case "set_fills": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId);
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      if (!("fills" in node)) throw new Error(`Node ${nodeId} does not support fills`);
      (node as any).fills = [makeSolidPaint(p.color, p.opacity != null ? p.opacity : undefined)];
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name },
      };
    }

    case "set_strokes": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId);
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      if (!("strokes" in node)) throw new Error(`Node ${nodeId} does not support strokes`);
      (node as any).strokes = [makeSolidPaint(p.color)];
      if (p.strokeWeight != null) (node as any).strokeWeight = p.strokeWeight;
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name },
      };
    }

    case "set_layout_properties": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId) as any;
      if (!node) throw new Error(`Node not found: ${nodeId}`);

      const applied: any = {};
      const setProp = (prop: string) => {
        if (!(prop in p)) return;
        if (!(prop in node)) throw new Error(`Node ${nodeId} does not support ${prop}`);
        node[prop] = p[prop];
        applied[prop] = p[prop];
      };

      [
        "layoutMode",
        "layoutWrap",
        "primaryAxisAlignItems",
        "counterAxisAlignItems",
        "primaryAxisSizingMode",
        "counterAxisSizingMode",
        "layoutSizingHorizontal",
        "layoutSizingVertical",
        "layoutAlign",
        "layoutGrow",
        "layoutPositioning",
        "paddingTop",
        "paddingRight",
        "paddingBottom",
        "paddingLeft",
        "itemSpacing",
        "counterAxisSpacing",
        "clipsContent",
      ].forEach(setProp);

      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name, applied },
      };
    }

    case "set_text_style": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId);
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      if (node.type !== "TEXT") throw new Error(`Node ${nodeId} is not a TEXT node`);

      await loadAllFonts(node);

      const currentFonts = node.characters.length > 0
        ? node.getRangeAllFontNames(0, node.characters.length)
        : [];
      const baseFont = typeof node.fontName !== "symbol"
        ? node.fontName
        : currentFonts[0] || { family: "Inter", style: "Regular" };

      const applied: any = {};

      if (p.fontFamily || p.fontStyle) {
        const nextFont = {
          family: p.fontFamily || baseFont.family,
          style: p.fontStyle || baseFont.style,
        };
        await figma.loadFontAsync(nextFont);
        node.fontName = nextFont;
        applied.fontFamily = nextFont.family;
        applied.fontStyle = nextFont.style;
      }
      if (p.fontSize != null) {
        node.fontSize = p.fontSize;
        applied.fontSize = p.fontSize;
      }
      if (p.textCase) {
        node.textCase = p.textCase;
        applied.textCase = p.textCase;
      }
      if (p.textAlignHorizontal) {
        node.textAlignHorizontal = p.textAlignHorizontal;
        applied.textAlignHorizontal = p.textAlignHorizontal;
      }
      if (p.textAlignVertical) {
        node.textAlignVertical = p.textAlignVertical;
        applied.textAlignVertical = p.textAlignVertical;
      }
      if (p.paragraphSpacing != null) {
        node.paragraphSpacing = p.paragraphSpacing;
        applied.paragraphSpacing = p.paragraphSpacing;
      }

      const lineHeight = parseLineHeight(p.lineHeight);
      if (lineHeight) {
        node.lineHeight = lineHeight;
        applied.lineHeight = lineHeight;
      }

      const letterSpacing = parseLetterSpacing(p.letterSpacing);
      if (letterSpacing) {
        node.letterSpacing = letterSpacing;
        applied.letterSpacing = letterSpacing;
      }

      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name, applied },
      };
    }

    case "set_effects": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId);
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      if (!("effects" in node)) throw new Error(`Node ${nodeId} does not support effects`);

      const effects = Array.isArray(p.effects) ? p.effects.map(makeEffect) : [];
      (node as any).effects = effects;

      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: {
          id: node.id,
          name: node.name,
          effectCount: effects.length,
          effectTypes: effects.map((effect: any) => effect.type),
        },
      };
    }

    case "apply_styles": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId) as any;
      if (!node) throw new Error(`Node not found: ${nodeId}`);

      const applied: any = {};

      if (p.fillStyleId) {
        if (typeof node.setFillStyleIdAsync !== "function") {
          throw new Error(`Node ${nodeId} does not support fill styles`);
        }
        await node.setFillStyleIdAsync(p.fillStyleId);
        applied.fillStyleId = p.fillStyleId;
      }
      if (p.strokeStyleId) {
        if (typeof node.setStrokeStyleIdAsync !== "function") {
          throw new Error(`Node ${nodeId} does not support stroke styles`);
        }
        await node.setStrokeStyleIdAsync(p.strokeStyleId);
        applied.strokeStyleId = p.strokeStyleId;
      }
      if (p.effectStyleId) {
        if (typeof node.setEffectStyleIdAsync !== "function") {
          throw new Error(`Node ${nodeId} does not support effect styles`);
        }
        await node.setEffectStyleIdAsync(p.effectStyleId);
        applied.effectStyleId = p.effectStyleId;
      }
      if (p.textStyleId) {
        if (node.type !== "TEXT") throw new Error(`Node ${nodeId} is not a TEXT node`);
        await loadAllFonts(node);
        await node.setTextStyleIdAsync(p.textStyleId);
        applied.textStyleId = p.textStyleId;
      }

      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name, applied },
      };
    }

    case "create_instance": {
      const p = request.params || {};
      const component = await resolveComponentNode({
        componentId: p.componentId,
        componentKey: p.componentKey,
        name: p.name,
      });
      const parent = await getParentNode(p.parentId);
      const instance = component.createInstance();

      if (parent !== figma.currentPage || instance.parent !== parent) {
        (parent as any).appendChild(instance);
      }
      if (p.x != null) instance.x = p.x;
      if (p.y != null) instance.y = p.y;

      if (p.componentProperties && typeof p.componentProperties === "object") {
        const mapped = mapComponentProperties(instance, p.componentProperties);
        if (Object.keys(mapped).length > 0) {
          instance.setProperties(mapped);
        }
      }

      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: {
          id: instance.id,
          name: instance.name,
          type: instance.type,
          mainComponentId: component.id,
          componentPropertyCount: Object.keys(instance.componentProperties || {}).length,
          bounds: getBounds(instance),
        },
      };
    }

    case "move_nodes": {
      const p = request.params || {};
      const nodeIds = request.nodeIds || [];
      if (nodeIds.length === 0) throw new Error("nodeIds is required");
      const results: any[] = [];
      for (const nid of nodeIds) {
        const n = await figma.getNodeByIdAsync(nid) as any;
        if (!n) { results.push({ nodeId: nid, error: "Node not found" }); continue; }
        if (!("x" in n)) { results.push({ nodeId: nid, error: "Node does not support position" }); continue; }
        if (p.x != null) n.x = p.x;
        if (p.y != null) n.y = p.y;
        results.push({ nodeId: nid, x: n.x, y: n.y });
      }
      figma.commitUndo();
      return { type: request.type, requestId: request.requestId, data: { results } };
    }

    case "resize_nodes": {
      const p = request.params || {};
      const nodeIds = request.nodeIds || [];
      if (nodeIds.length === 0) throw new Error("nodeIds is required");
      const results: any[] = [];
      for (const nid of nodeIds) {
        const n = await figma.getNodeByIdAsync(nid) as any;
        if (!n) { results.push({ nodeId: nid, error: "Node not found" }); continue; }
        if (!("resize" in n)) { results.push({ nodeId: nid, error: "Node does not support resize" }); continue; }
        const w = p.width != null ? p.width : n.width;
        const h = p.height != null ? p.height : n.height;
        n.resize(w, h);
        results.push({ nodeId: nid, width: n.width, height: n.height });
      }
      figma.commitUndo();
      return { type: request.type, requestId: request.requestId, data: { results } };
    }

    case "delete_nodes": {
      const nodeIds = request.nodeIds || [];
      if (nodeIds.length === 0) throw new Error("nodeIds is required");
      const results: any[] = [];
      for (const nid of nodeIds) {
        const n = await figma.getNodeByIdAsync(nid);
        if (!n) { results.push({ nodeId: nid, error: "Node not found" }); continue; }
        n.remove();
        results.push({ nodeId: nid, deleted: true });
      }
      figma.commitUndo();
      return { type: request.type, requestId: request.requestId, data: { results } };
    }

    case "rename_node": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId);
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      node.name = p.name;
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: node.id, name: node.name },
      };
    }

    case "clone_node": {
      const p = request.params || {};
      const nodeId = request.nodeIds && request.nodeIds[0];
      if (!nodeId) throw new Error("nodeId is required");
      const node = await figma.getNodeByIdAsync(nodeId) as any;
      if (!node) throw new Error(`Node not found: ${nodeId}`);
      const clone = node.clone();
      if (p.x != null) clone.x = p.x;
      if (p.y != null) clone.y = p.y;
      if (p.parentId) {
        const parent = await getParentNode(p.parentId);
        (parent as any).appendChild(clone);
      }
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: clone.id, name: clone.name, type: clone.type, bounds: getBounds(clone) },
      };
    }

    case "import_image": {
      const p = request.params || {};
      if (!p.imageData) throw new Error("imageData (base64) is required");
      const parent = await getParentNode(p.parentId);
      const bytes = base64ToBytes(p.imageData);
      const image = figma.createImage(bytes);
      const rect = figma.createRectangle();
      rect.resize(p.width || 200, p.height || 200);
      rect.x = p.x != null ? p.x : 0;
      rect.y = p.y != null ? p.y : 0;
      if (p.name) rect.name = p.name;
      rect.fills = [{ type: "IMAGE", imageHash: image.hash, scaleMode: p.scaleMode || "FILL" }];
      (parent as any).appendChild(rect);
      figma.commitUndo();
      return {
        type: request.type,
        requestId: request.requestId,
        data: { id: rect.id, name: rect.name, type: rect.type, bounds: getBounds(rect) },
      };
    }

    default:
      return null;
  }
};
