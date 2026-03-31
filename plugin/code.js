var __async = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
(function() {
  "use strict";
  const isMixed = (value) => typeof value === "symbol";
  const toHex = (color) => {
    const clamp = (v) => Math.min(255, Math.max(0, Math.round(v * 255)));
    const [r, g, b] = [clamp(color.r), clamp(color.g), clamp(color.b)];
    return `#${[r, g, b].map((v) => v.toString(16).padStart(2, "0")).join("")}`;
  };
  const serializePaints = (paints) => {
    if (isMixed(paints)) return "mixed";
    if (!paints || !Array.isArray(paints)) return void 0;
    const result = paints.filter((paint) => paint.type === "SOLID" && "color" in paint).map((paint) => {
      const hex = toHex(paint.color);
      const opacity = paint.opacity != null ? paint.opacity : 1;
      if (opacity === 1) return hex;
      return hex + Math.round(opacity * 255).toString(16).padStart(2, "0");
    });
    return result.length > 0 ? result : void 0;
  };
  const getBounds = (node) => {
    if ("x" in node && "y" in node && "width" in node && "height" in node) {
      return { x: node.x, y: node.y, width: node.width, height: node.height };
    }
    return void 0;
  };
  const serializeStyles = (node) => {
    const styles = {};
    if ("fills" in node) {
      const fills = serializePaints(node.fills);
      if (fills !== void 0) styles.fills = fills;
    }
    if ("strokes" in node) {
      const strokes = serializePaints(node.strokes);
      if (strokes !== void 0) styles.strokes = strokes;
    }
    if ("cornerRadius" in node) {
      const cr = isMixed(node.cornerRadius) ? "mixed" : node.cornerRadius;
      if (cr !== 0) styles.cornerRadius = cr;
    }
    if ("paddingLeft" in node) {
      styles.padding = {
        top: node.paddingTop,
        right: node.paddingRight,
        bottom: node.paddingBottom,
        left: node.paddingLeft
      };
    }
    return styles;
  };
  const serializeLineHeight = (lineHeight) => {
    if (isMixed(lineHeight)) return "mixed";
    if (!lineHeight || lineHeight.unit === "AUTO") return void 0;
    return { value: lineHeight.value, unit: lineHeight.unit };
  };
  const serializeLetterSpacing = (letterSpacing) => {
    if (isMixed(letterSpacing)) return "mixed";
    if (!letterSpacing || letterSpacing.value === 0) return void 0;
    return { value: letterSpacing.value, unit: letterSpacing.unit };
  };
  const serializeText = (node, base) => {
    let fontFamily;
    let fontStyle;
    if (typeof node.fontName === "symbol") {
      fontFamily = "mixed";
      fontStyle = "mixed";
    } else if (node.fontName) {
      fontFamily = node.fontName.family;
      fontStyle = node.fontName.style;
    }
    return Object.assign({}, base, {
      characters: node.characters,
      styles: Object.assign({}, base.styles, {
        fontSize: isMixed(node.fontSize) ? "mixed" : node.fontSize,
        fontFamily,
        fontStyle,
        fontWeight: isMixed(node.fontWeight) ? "mixed" : node.fontWeight,
        textDecoration: isMixed(node.textDecoration) ? "mixed" : node.textDecoration !== "NONE" ? node.textDecoration : void 0,
        lineHeight: serializeLineHeight(node.lineHeight),
        letterSpacing: serializeLetterSpacing(node.letterSpacing),
        textAlignHorizontal: isMixed(node.textAlignHorizontal) ? "mixed" : node.textAlignHorizontal
      })
    });
  };
  const serializeNode = (node) => {
    const base = {
      id: node.id,
      name: node.name,
      type: node.type,
      bounds: getBounds(node),
      styles: serializeStyles(node)
    };
    if (node.type === "TEXT") return serializeText(node, base);
    if ("children" in node) {
      return Object.assign({}, base, {
        children: node.children.map((child) => serializeNode(child))
      });
    }
    return base;
  };
  const serializeVariableValue = (value) => {
    if (typeof value !== "object" || value === null) return value;
    if ("type" in value && value.type === "VARIABLE_ALIAS") {
      return { type: "VARIABLE_ALIAS", id: value.id };
    }
    if ("r" in value && "g" in value && "b" in value) {
      return {
        type: "COLOR",
        r: value.r,
        g: value.g,
        b: value.b,
        a: "a" in value ? value.a : 1
      };
    }
    return value;
  };
  const handleReadRequest = (request) => __async(null, null, function* () {
    switch (request.type) {
      case "get_document":
        return {
          type: request.type,
          requestId: request.requestId,
          data: serializeNode(figma.currentPage)
        };
      case "get_selection":
        return {
          type: request.type,
          requestId: request.requestId,
          data: figma.currentPage.selection.map((node) => serializeNode(node))
        };
      case "get_node": {
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeIds is required for get_node");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node || node.type === "DOCUMENT")
          throw new Error(`Node not found: ${nodeId}`);
        return {
          type: request.type,
          requestId: request.requestId,
          data: serializeNode(node)
        };
      }
      case "get_styles": {
        const [paintStyles, textStyles, effectStyles, gridStyles] = yield Promise.all([
          figma.getLocalPaintStylesAsync(),
          figma.getLocalTextStylesAsync(),
          figma.getLocalEffectStylesAsync(),
          figma.getLocalGridStylesAsync()
        ]);
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            paints: paintStyles.map((s) => ({
              id: s.id,
              name: s.name,
              paints: s.paints
            })),
            text: textStyles.map((s) => ({
              id: s.id,
              name: s.name,
              fontSize: s.fontSize,
              fontFamily: s.fontName ? s.fontName.family : void 0,
              fontStyle: s.fontName ? s.fontName.style : void 0,
              textDecoration: s.textDecoration !== "NONE" ? s.textDecoration : void 0,
              lineHeight: s.lineHeight,
              letterSpacing: s.letterSpacing
            })),
            effects: effectStyles.map((s) => ({
              id: s.id,
              name: s.name,
              effects: s.effects
            })),
            grids: gridStyles.map((s) => ({
              id: s.id,
              name: s.name,
              layoutGrids: s.layoutGrids
            }))
          }
        };
      }
      case "get_metadata":
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            fileName: figma.root.name,
            currentPageId: figma.currentPage.id,
            currentPageName: figma.currentPage.name,
            pageCount: figma.root.children.length,
            pages: figma.root.children.map((page) => ({
              id: page.id,
              name: page.name
            }))
          }
        };
      case "get_design_context": {
        const depth = request.params && request.params.depth != null ? request.params.depth : 2;
        const detail = request.params && request.params.detail || "full";
        const serializeForDetail = (n) => {
          const base = { id: n.id, name: n.name, type: n.type, bounds: getBounds(n) };
          if (detail === "minimal") return base;
          const styles = serializeStyles(n);
          const result = Object.assign({}, base);
          if (Object.keys(styles).length > 0) result.styles = styles;
          if ("opacity" in n && n.opacity !== 1) result.opacity = n.opacity;
          if ("visible" in n && !n.visible) result.visible = false;
          if (detail === "compact") return result;
          return serializeNode(n);
        };
        const serializeWithDepth = (node, currentDepth) => __async(null, null, function* () {
          if (detail === "full") {
            const serialized2 = serializeNode(node);
            if (currentDepth >= depth && serialized2.children) {
              return Object.assign({}, serialized2, {
                children: void 0,
                childCount: node.children ? node.children.length : 0
              });
            }
            if (serialized2.children) {
              const childNodes = yield Promise.all(
                serialized2.children.map(
                  (child) => figma.getNodeByIdAsync(child.id)
                )
              );
              const serializedChildren2 = yield Promise.all(
                childNodes.filter((n) => n !== null && n.type !== "DOCUMENT").map((n) => serializeWithDepth(n, currentDepth + 1))
              );
              return Object.assign({}, serialized2, { children: serializedChildren2 });
            }
            return serialized2;
          }
          const serialized = serializeForDetail(node);
          const hasChildren = "children" in node && node.children.length > 0;
          if (!hasChildren) return serialized;
          if (currentDepth >= depth) {
            return Object.assign({}, serialized, { childCount: node.children.length });
          }
          const serializedChildren = yield Promise.all(
            node.children.filter((n) => n.type !== "DOCUMENT").map((n) => serializeWithDepth(n, currentDepth + 1))
          );
          return Object.assign({}, serialized, { children: serializedChildren });
        });
        const selection = figma.currentPage.selection;
        const contextNodes = selection.length > 0 ? yield Promise.all(
          selection.map((node) => serializeWithDepth(node, 0))
        ) : [yield serializeWithDepth(figma.currentPage, 0)];
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            fileName: figma.root.name,
            currentPage: {
              id: figma.currentPage.id,
              name: figma.currentPage.name
            },
            selectionCount: selection.length,
            context: contextNodes
          }
        };
      }
      case "get_variable_defs": {
        const collections = yield figma.variables.getLocalVariableCollectionsAsync();
        const variableData = yield Promise.all(
          collections.map((collection) => __async(null, null, function* () {
            const variables = yield Promise.all(
              collection.variableIds.map(
                (id) => figma.variables.getVariableByIdAsync(id)
              )
            );
            return {
              id: collection.id,
              name: collection.name,
              modes: collection.modes.map((mode) => ({
                modeId: mode.modeId,
                name: mode.name
              })),
              variables: variables.filter((v) => v !== null).map((variable) => ({
                id: variable.id,
                name: variable.name,
                resolvedType: variable.resolvedType,
                valuesByMode: Object.fromEntries(
                  Object.entries(variable.valuesByMode).map(
                    ([modeId, value]) => [
                      modeId,
                      serializeVariableValue(value)
                    ]
                  )
                )
              }))
            };
          }))
        );
        return {
          type: request.type,
          requestId: request.requestId,
          data: { collections: variableData }
        };
      }
      case "get_screenshot": {
        const format = request.params && request.params.format ? request.params.format : "PNG";
        const scale = request.params && request.params.scale != null ? request.params.scale : 2;
        let targetNodes;
        if (request.nodeIds && request.nodeIds.length > 0) {
          const nodes = yield Promise.all(
            request.nodeIds.map((id) => figma.getNodeByIdAsync(id))
          );
          targetNodes = nodes.filter(
            (n) => n !== null && n.type !== "DOCUMENT" && n.type !== "PAGE"
          );
        } else {
          targetNodes = figma.currentPage.selection.slice();
        }
        if (targetNodes.length === 0)
          throw new Error(
            "No nodes to export. Select nodes or provide nodeIds."
          );
        const exports$1 = yield Promise.all(
          targetNodes.map((node) => __async(null, null, function* () {
            const settings = format === "SVG" ? { format: "SVG" } : format === "PDF" ? { format: "PDF" } : format === "JPG" ? {
              format: "JPG",
              constraint: { type: "SCALE", value: scale }
            } : {
              format: "PNG",
              constraint: { type: "SCALE", value: scale }
            };
            const bytes = yield node.exportAsync(settings);
            const base64 = figma.base64Encode(bytes);
            return {
              nodeId: node.id,
              nodeName: node.name,
              format,
              base64,
              width: node.width,
              height: node.height
            };
          }))
        );
        return {
          type: request.type,
          requestId: request.requestId,
          data: { exports: exports$1 }
        };
      }
      case "get_nodes_info": {
        if (!request.nodeIds || request.nodeIds.length === 0)
          throw new Error("nodeIds is required for get_nodes_info");
        const nodes = yield Promise.all(
          request.nodeIds.map((id) => figma.getNodeByIdAsync(id))
        );
        return {
          type: request.type,
          requestId: request.requestId,
          data: nodes.filter((n) => n !== null && n.type !== "DOCUMENT").map((n) => serializeNode(n))
        };
      }
      case "get_local_components": {
        const pages = figma.root.children;
        const allComponents = [];
        const componentSetsMap = /* @__PURE__ */ new Map();
        for (let i = 0; i < pages.length; i++) {
          const page = pages[i];
          yield page.loadAsync();
          const pageNodes = page.findAllWithCriteria({
            types: ["COMPONENT", "COMPONENT_SET"]
          });
          for (const n of pageNodes) {
            if (n.type === "COMPONENT_SET") {
              componentSetsMap.set(n.id, {
                id: n.id,
                name: n.name,
                key: "key" in n ? n.key : null
              });
            } else {
              const parentIsSet = n.parent && n.parent.type === "COMPONENT_SET";
              allComponents.push({
                id: n.id,
                name: n.name,
                key: "key" in n ? n.key : null,
                componentSetId: parentIsSet ? n.parent.id : null,
                variantProperties: "variantProperties" in n ? n.variantProperties : null
              });
            }
          }
          figma.ui.postMessage({
            type: "progress_update",
            requestId: request.requestId,
            progress: Math.round((i + 1) / pages.length * 90) + 1,
            message: `Scanned ${page.name}: ${allComponents.length} components so far`
          });
          yield new Promise((r) => setTimeout(r, 0));
        }
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            count: allComponents.length,
            components: allComponents,
            componentSets: Array.from(componentSetsMap.values())
          }
        };
      }
      case "get_pages":
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            currentPageId: figma.currentPage.id,
            pages: figma.root.children.map((page) => ({
              id: page.id,
              name: page.name
            }))
          }
        };
      case "get_viewport":
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            center: { x: figma.viewport.center.x, y: figma.viewport.center.y },
            zoom: figma.viewport.zoom,
            bounds: {
              x: figma.viewport.bounds.x,
              y: figma.viewport.bounds.y,
              width: figma.viewport.bounds.width,
              height: figma.viewport.bounds.height
            }
          }
        };
      case "get_fonts": {
        const fontMap = /* @__PURE__ */ new Map();
        const collectFonts = (n) => {
          if (n.type === "TEXT") {
            const fontName = n.fontName;
            if (typeof fontName !== "symbol" && fontName) {
              const key = `${fontName.family}::${fontName.style}`;
              if (!fontMap.has(key)) {
                fontMap.set(key, { family: fontName.family, style: fontName.style, nodeCount: 0 });
              }
              fontMap.get(key).nodeCount++;
            }
          }
          if ("children" in n) n.children.forEach(collectFonts);
        };
        collectFonts(figma.currentPage);
        const fonts = Array.from(fontMap.values()).sort((a, b) => b.nodeCount - a.nodeCount);
        return {
          type: request.type,
          requestId: request.requestId,
          data: { count: fonts.length, fonts }
        };
      }
      case "search_nodes": {
        const query = request.params && request.params.query ? request.params.query.toLowerCase() : "";
        const scopeNodeId = request.params && request.params.nodeId;
        const types = request.params && request.params.types ? request.params.types : [];
        const limit = request.params && request.params.limit ? request.params.limit : 50;
        const root = scopeNodeId ? yield figma.getNodeByIdAsync(scopeNodeId) : figma.currentPage;
        if (!root) throw new Error(`Node not found: ${scopeNodeId}`);
        const results = [];
        const search = (n) => __async(null, null, function* () {
          if (results.length >= limit) return;
          if (n !== root) {
            const nameMatch = !query || n.name.toLowerCase().includes(query);
            const typeMatch = types.length === 0 || types.includes(n.type);
            if (nameMatch && typeMatch) {
              results.push({
                id: n.id,
                name: n.name,
                type: n.type,
                bounds: getBounds(n)
              });
            }
          }
          if (results.length < limit && "children" in n) {
            for (const child of n.children) yield search(child);
          }
        });
        yield search(root);
        return {
          type: request.type,
          requestId: request.requestId,
          data: { count: results.length, nodes: results }
        };
      }
      case "get_reactions": {
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required for get_reactions");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node || node.type === "DOCUMENT") throw new Error(`Node not found: ${nodeId}`);
        const reactions = "reactions" in node ? node.reactions : [];
        return {
          type: request.type,
          requestId: request.requestId,
          data: { nodeId: node.id, name: node.name, reactions }
        };
      }
      case "get_annotations": {
        const nodeId = request.params && request.params.nodeId;
        const nodeAnnotations = (n) => {
          const anns = n.annotations;
          return Array.isArray(anns) ? anns : null;
        };
        if (nodeId) {
          const node = yield figma.getNodeByIdAsync(nodeId);
          if (!node) throw new Error(`Node not found: ${nodeId}`);
          const mergedAnnotations = [];
          const collect = (n) => __async(null, null, function* () {
            const anns = nodeAnnotations(n);
            if (anns)
              for (const a of anns)
                mergedAnnotations.push({ nodeId: n.id, annotation: a });
            if ("children" in n)
              for (const child of n.children) yield collect(child);
          });
          yield collect(node);
          return {
            type: request.type,
            requestId: request.requestId,
            data: {
              nodeId: node.id,
              name: node.name,
              annotations: mergedAnnotations
            }
          };
        }
        const annotated = [];
        const processNode = (n) => __async(null, null, function* () {
          const anns = nodeAnnotations(n);
          if (anns && anns.length > 0)
            annotated.push({ nodeId: n.id, name: n.name, annotations: anns });
          if ("children" in n)
            for (const child of n.children) yield processNode(child);
        });
        yield processNode(figma.currentPage);
        return {
          type: request.type,
          requestId: request.requestId,
          data: { annotatedNodes: annotated }
        };
      }
      case "scan_text_nodes": {
        const nodeId = request.params && request.params.nodeId;
        if (!nodeId) throw new Error("nodeId is required for scan_text_nodes");
        const root = yield figma.getNodeByIdAsync(nodeId);
        if (!root) throw new Error(`Node not found: ${nodeId}`);
        const textNodes = [];
        const findText = (n) => __async(null, null, function* () {
          if (n.type === "TEXT") {
            textNodes.push({
              id: n.id,
              name: n.name,
              characters: n.characters,
              fontSize: n.fontSize,
              fontName: n.fontName
            });
          }
          if ("children" in n)
            for (const child of n.children) yield findText(child);
        });
        figma.ui.postMessage({
          type: "progress_update",
          requestId: request.requestId,
          progress: 10,
          message: "Scanning text nodes..."
        });
        yield new Promise((r) => setTimeout(r, 0));
        yield findText(root);
        return {
          type: request.type,
          requestId: request.requestId,
          data: { count: textNodes.length, textNodes }
        };
      }
      case "scan_nodes_by_types": {
        const nodeId = request.params && request.params.nodeId;
        const types = request.params && request.params.types ? request.params.types : [];
        if (!nodeId)
          throw new Error("nodeId is required for scan_nodes_by_types");
        if (types.length === 0)
          throw new Error("types must be a non-empty array");
        const root = yield figma.getNodeByIdAsync(nodeId);
        if (!root) throw new Error(`Node not found: ${nodeId}`);
        const matchingNodes = [];
        const findByTypes = (n) => __async(null, null, function* () {
          if ("visible" in n && !n.visible) return;
          if (types.includes(n.type)) {
            matchingNodes.push({
              id: n.id,
              name: n.name,
              type: n.type,
              bbox: {
                x: "x" in n ? n.x : 0,
                y: "y" in n ? n.y : 0,
                width: "width" in n ? n.width : 0,
                height: "height" in n ? n.height : 0
              }
            });
          }
          if ("children" in n)
            for (const child of n.children) yield findByTypes(child);
        });
        figma.ui.postMessage({
          type: "progress_update",
          requestId: request.requestId,
          progress: 10,
          message: `Scanning for types: ${types.join(", ")}...`
        });
        yield new Promise((r) => setTimeout(r, 0));
        yield findByTypes(root);
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            count: matchingNodes.length,
            matchingNodes,
            searchedTypes: types
          }
        };
      }
      default:
        return null;
    }
  });
  const hexToRgb = (hex) => {
    const clean = hex.replace("#", "");
    return {
      r: parseInt(clean.slice(0, 2), 16) / 255,
      g: parseInt(clean.slice(2, 4), 16) / 255,
      b: parseInt(clean.slice(4, 6), 16) / 255,
      a: clean.length >= 8 ? parseInt(clean.slice(6, 8), 16) / 255 : 1
    };
  };
  const makeSolidPaint = (colorInput, opacityOverride) => {
    const { r, g, b, a } = typeof colorInput === "string" ? hexToRgb(colorInput) : { r: colorInput.r, g: colorInput.g, b: colorInput.b, a: colorInput.a != null ? colorInput.a : 1 };
    const eff = opacityOverride != null ? opacityOverride : a;
    const paint = { type: "SOLID", color: { r, g, b } };
    if (eff !== 1) paint.opacity = eff;
    return paint;
  };
  const getParentNode = (parentId) => __async(null, null, function* () {
    if (!parentId) return figma.currentPage;
    const parent = yield figma.getNodeByIdAsync(parentId);
    if (!parent) throw new Error(`Parent node not found: ${parentId}`);
    if (!("appendChild" in parent)) throw new Error(`Node ${parentId} cannot have children`);
    return parent;
  });
  const base64ToBytes = (b64) => {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    const lookup = {};
    for (let i = 0; i < chars.length; i++) lookup[chars[i]] = i;
    const clean = b64.replace(/[^A-Za-z0-9+/]/g, "");
    let outLen = Math.floor(clean.length * 3 / 4);
    if (b64.endsWith("==")) outLen -= 2;
    else if (b64.endsWith("=")) outLen -= 1;
    const bytes = new Uint8Array(outLen);
    let j = 0;
    for (let i = 0; i < clean.length; i += 4) {
      const a = lookup[clean[i]] || 0;
      const bv = lookup[clean[i + 1]] || 0;
      const c = lookup[clean[i + 2]] || 0;
      const d = lookup[clean[i + 3]] || 0;
      bytes[j++] = a << 2 | bv >> 4;
      if (j < outLen) bytes[j++] = (bv & 15) << 4 | c >> 2;
      if (j < outLen) bytes[j++] = (c & 3) << 6 | d;
    }
    return bytes;
  };
  const makeRGBA = (colorInput, opacityOverride) => {
    const { r, g, b, a } = typeof colorInput === "string" ? hexToRgb(colorInput) : { r: colorInput.r, g: colorInput.g, b: colorInput.b, a: colorInput.a != null ? colorInput.a : 1 };
    return { r, g, b, a: opacityOverride != null ? opacityOverride : a };
  };
  const loadAllFonts = (node) => __async(null, null, function* () {
    if (node.hasMissingFont) {
      throw new Error(`Text node ${node.id} has missing fonts`);
    }
    if (typeof node.fontName !== "symbol" && node.fontName) {
      yield figma.loadFontAsync(node.fontName);
      return;
    }
    const seen = /* @__PURE__ */ new Set();
    const fonts = node.getRangeAllFontNames(0, node.characters.length);
    for (const font of fonts) {
      const key = `${font.family}::${font.style}`;
      if (seen.has(key)) continue;
      seen.add(key);
      yield figma.loadFontAsync(font);
    }
  });
  const parseLineHeight = (value) => {
    if (!value || typeof value !== "object") return void 0;
    if (value.unit === "AUTO") return { unit: "AUTO" };
    if (!value.unit || value.value == null) return void 0;
    return { unit: value.unit, value: value.value };
  };
  const parseLetterSpacing = (value) => {
    if (!value || typeof value !== "object") return void 0;
    if (!value.unit || value.value == null) return void 0;
    return { unit: value.unit, value: value.value };
  };
  const makeEffect = (effect) => {
    switch (effect.type) {
      case "DROP_SHADOW":
        return {
          type: "DROP_SHADOW",
          color: makeRGBA(effect.color),
          offset: { x: effect.x != null ? effect.x : 0, y: effect.y != null ? effect.y : 4 },
          radius: effect.radius != null ? effect.radius : 16,
          spread: effect.spread != null ? effect.spread : 0,
          visible: effect.visible != null ? effect.visible : true,
          blendMode: effect.blendMode || "NORMAL",
          showShadowBehindNode: effect.showShadowBehindNode === true
        };
      case "INNER_SHADOW":
        return {
          type: "INNER_SHADOW",
          color: makeRGBA(effect.color),
          offset: { x: effect.x != null ? effect.x : 0, y: effect.y != null ? effect.y : 2 },
          radius: effect.radius != null ? effect.radius : 12,
          spread: effect.spread != null ? effect.spread : 0,
          visible: effect.visible != null ? effect.visible : true,
          blendMode: effect.blendMode || "NORMAL"
        };
      case "LAYER_BLUR":
      case "BACKGROUND_BLUR":
        return {
          type: effect.type,
          radius: effect.radius != null ? effect.radius : 12,
          visible: effect.visible != null ? effect.visible : true,
          blurType: "NORMAL"
        };
      default:
        throw new Error(`Unsupported effect type: ${effect.type}`);
    }
  };
  const toComponentNode = (node) => {
    if (!node) return null;
    if (node.type === "COMPONENT") return node;
    if (node.type === "COMPONENT_SET") return node.defaultVariant;
    return null;
  };
  const resolveComponentNode = (selector) => __async(null, null, function* () {
    if (selector.componentId) {
      const node = yield figma.getNodeByIdAsync(selector.componentId);
      const component = toComponentNode(node);
      if (component) return component;
      throw new Error(`Component not found: ${selector.componentId}`);
    }
    const pages = figma.root.children;
    let caseInsensitiveMatch = null;
    for (const page of pages) {
      yield page.loadAsync();
      const candidates = page.findAllWithCriteria({
        types: ["COMPONENT", "COMPONENT_SET"]
      });
      for (const candidate of candidates) {
        const component = toComponentNode(candidate);
        if (!component) continue;
        if (selector.componentKey && "key" in component && component.key === selector.componentKey) {
          return component;
        }
        if (selector.name) {
          if (component.name === selector.name) return component;
          if (!caseInsensitiveMatch && component.name.toLowerCase() === selector.name.toLowerCase()) {
            caseInsensitiveMatch = component;
          }
        }
      }
    }
    if (caseInsensitiveMatch) return caseInsensitiveMatch;
    if (selector.componentKey) {
      throw new Error(`Component key not found: ${selector.componentKey}`);
    }
    throw new Error(`Component name not found: ${selector.name}`);
  });
  const mapComponentProperties = (instance, rawProps) => {
    const mapped = {};
    const available = instance.componentProperties || {};
    const availableKeys = Object.keys(available);
    for (const [rawKey, rawValue] of Object.entries(rawProps)) {
      if (typeof rawValue !== "string" && typeof rawValue !== "boolean") {
        throw new Error(`Component property ${rawKey} must be a string or boolean`);
      }
      let resolvedKey = availableKeys.find((key) => key === rawKey);
      if (!resolvedKey) {
        resolvedKey = availableKeys.find((key) => key.split("#")[0] === rawKey);
      }
      if (!resolvedKey) {
        throw new Error(`Component property not found on instance: ${rawKey}`);
      }
      mapped[resolvedKey] = rawValue;
    }
    return mapped;
  };
  const handleWriteRequest = (request) => __async(null, null, function* () {
    switch (request.type) {
      case "create_frame": {
        const p = request.params || {};
        const parent = yield getParentNode(p.parentId);
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
        parent.appendChild(frame);
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: frame.id, name: frame.name, type: frame.type, bounds: getBounds(frame) }
        };
      }
      case "create_rectangle": {
        const p = request.params || {};
        const parent = yield getParentNode(p.parentId);
        const rect = figma.createRectangle();
        rect.resize(p.width || 100, p.height || 100);
        rect.x = p.x != null ? p.x : 0;
        rect.y = p.y != null ? p.y : 0;
        if (p.name) rect.name = p.name;
        if (p.fillColor) rect.fills = [makeSolidPaint(p.fillColor)];
        if (p.cornerRadius != null) rect.cornerRadius = p.cornerRadius;
        parent.appendChild(rect);
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: rect.id, name: rect.name, type: rect.type, bounds: getBounds(rect) }
        };
      }
      case "create_ellipse": {
        const p = request.params || {};
        const parent = yield getParentNode(p.parentId);
        const ellipse = figma.createEllipse();
        ellipse.resize(p.width || 100, p.height || 100);
        ellipse.x = p.x != null ? p.x : 0;
        ellipse.y = p.y != null ? p.y : 0;
        if (p.name) ellipse.name = p.name;
        if (p.fillColor) ellipse.fills = [makeSolidPaint(p.fillColor)];
        parent.appendChild(ellipse);
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: ellipse.id, name: ellipse.name, type: ellipse.type, bounds: getBounds(ellipse) }
        };
      }
      case "create_text": {
        const p = request.params || {};
        const parent = yield getParentNode(p.parentId);
        const fontFamily = p.fontFamily || "Inter";
        const fontStyle = p.fontStyle || "Regular";
        yield figma.loadFontAsync({ family: fontFamily, style: fontStyle });
        const textNode = figma.createText();
        textNode.fontName = { family: fontFamily, style: fontStyle };
        if (p.fontSize) textNode.fontSize = p.fontSize;
        textNode.characters = p.text || "";
        textNode.x = p.x != null ? p.x : 0;
        textNode.y = p.y != null ? p.y : 0;
        if (p.name) textNode.name = p.name;
        if (p.fillColor) textNode.fills = [makeSolidPaint(p.fillColor)];
        parent.appendChild(textNode);
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: textNode.id, name: textNode.name, type: textNode.type, bounds: getBounds(textNode) }
        };
      }
      case "set_text": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        if (node.type !== "TEXT") throw new Error(`Node ${nodeId} is not a TEXT node`);
        const fontName = typeof node.fontName === "symbol" ? { family: "Inter", style: "Regular" } : node.fontName;
        yield figma.loadFontAsync(fontName);
        node.characters = p.text;
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: node.id, name: node.name, characters: node.characters }
        };
      }
      case "set_fills": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        if (!("fills" in node)) throw new Error(`Node ${nodeId} does not support fills`);
        node.fills = [makeSolidPaint(p.color, p.opacity != null ? p.opacity : void 0)];
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: node.id, name: node.name }
        };
      }
      case "set_strokes": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        if (!("strokes" in node)) throw new Error(`Node ${nodeId} does not support strokes`);
        node.strokes = [makeSolidPaint(p.color)];
        if (p.strokeWeight != null) node.strokeWeight = p.strokeWeight;
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: node.id, name: node.name }
        };
      }
      case "set_layout_properties": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        const applied = {};
        const setProp = (prop) => {
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
          "clipsContent"
        ].forEach(setProp);
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: node.id, name: node.name, applied }
        };
      }
      case "set_text_style": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        if (node.type !== "TEXT") throw new Error(`Node ${nodeId} is not a TEXT node`);
        yield loadAllFonts(node);
        const currentFonts = node.characters.length > 0 ? node.getRangeAllFontNames(0, node.characters.length) : [];
        const baseFont = typeof node.fontName !== "symbol" ? node.fontName : currentFonts[0] || { family: "Inter", style: "Regular" };
        const applied = {};
        if (p.fontFamily || p.fontStyle) {
          const nextFont = {
            family: p.fontFamily || baseFont.family,
            style: p.fontStyle || baseFont.style
          };
          yield figma.loadFontAsync(nextFont);
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
          data: { id: node.id, name: node.name, applied }
        };
      }
      case "set_effects": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        if (!("effects" in node)) throw new Error(`Node ${nodeId} does not support effects`);
        const effects = Array.isArray(p.effects) ? p.effects.map(makeEffect) : [];
        node.effects = effects;
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: {
            id: node.id,
            name: node.name,
            effectCount: effects.length,
            effectTypes: effects.map((effect) => effect.type)
          }
        };
      }
      case "apply_styles": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        const applied = {};
        if (p.fillStyleId) {
          if (typeof node.setFillStyleIdAsync !== "function") {
            throw new Error(`Node ${nodeId} does not support fill styles`);
          }
          yield node.setFillStyleIdAsync(p.fillStyleId);
          applied.fillStyleId = p.fillStyleId;
        }
        if (p.strokeStyleId) {
          if (typeof node.setStrokeStyleIdAsync !== "function") {
            throw new Error(`Node ${nodeId} does not support stroke styles`);
          }
          yield node.setStrokeStyleIdAsync(p.strokeStyleId);
          applied.strokeStyleId = p.strokeStyleId;
        }
        if (p.effectStyleId) {
          if (typeof node.setEffectStyleIdAsync !== "function") {
            throw new Error(`Node ${nodeId} does not support effect styles`);
          }
          yield node.setEffectStyleIdAsync(p.effectStyleId);
          applied.effectStyleId = p.effectStyleId;
        }
        if (p.textStyleId) {
          if (node.type !== "TEXT") throw new Error(`Node ${nodeId} is not a TEXT node`);
          yield loadAllFonts(node);
          yield node.setTextStyleIdAsync(p.textStyleId);
          applied.textStyleId = p.textStyleId;
        }
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: node.id, name: node.name, applied }
        };
      }
      case "create_instance": {
        const p = request.params || {};
        const component = yield resolveComponentNode({
          componentId: p.componentId,
          componentKey: p.componentKey,
          name: p.name
        });
        const parent = yield getParentNode(p.parentId);
        const instance = component.createInstance();
        if (parent !== figma.currentPage || instance.parent !== parent) {
          parent.appendChild(instance);
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
            bounds: getBounds(instance)
          }
        };
      }
      case "move_nodes": {
        const p = request.params || {};
        const nodeIds = request.nodeIds || [];
        if (nodeIds.length === 0) throw new Error("nodeIds is required");
        const results = [];
        for (const nid of nodeIds) {
          const n = yield figma.getNodeByIdAsync(nid);
          if (!n) {
            results.push({ nodeId: nid, error: "Node not found" });
            continue;
          }
          if (!("x" in n)) {
            results.push({ nodeId: nid, error: "Node does not support position" });
            continue;
          }
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
        const results = [];
        for (const nid of nodeIds) {
          const n = yield figma.getNodeByIdAsync(nid);
          if (!n) {
            results.push({ nodeId: nid, error: "Node not found" });
            continue;
          }
          if (!("resize" in n)) {
            results.push({ nodeId: nid, error: "Node does not support resize" });
            continue;
          }
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
        const results = [];
        for (const nid of nodeIds) {
          const n = yield figma.getNodeByIdAsync(nid);
          if (!n) {
            results.push({ nodeId: nid, error: "Node not found" });
            continue;
          }
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
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        node.name = p.name;
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: node.id, name: node.name }
        };
      }
      case "clone_node": {
        const p = request.params || {};
        const nodeId = request.nodeIds && request.nodeIds[0];
        if (!nodeId) throw new Error("nodeId is required");
        const node = yield figma.getNodeByIdAsync(nodeId);
        if (!node) throw new Error(`Node not found: ${nodeId}`);
        const clone = node.clone();
        if (p.x != null) clone.x = p.x;
        if (p.y != null) clone.y = p.y;
        if (p.parentId) {
          const parent = yield getParentNode(p.parentId);
          parent.appendChild(clone);
        }
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: clone.id, name: clone.name, type: clone.type, bounds: getBounds(clone) }
        };
      }
      case "import_image": {
        const p = request.params || {};
        if (!p.imageData) throw new Error("imageData (base64) is required");
        const parent = yield getParentNode(p.parentId);
        const bytes = base64ToBytes(p.imageData);
        const image = figma.createImage(bytes);
        const rect = figma.createRectangle();
        rect.resize(p.width || 200, p.height || 200);
        rect.x = p.x != null ? p.x : 0;
        rect.y = p.y != null ? p.y : 0;
        if (p.name) rect.name = p.name;
        rect.fills = [{ type: "IMAGE", imageHash: image.hash, scaleMode: p.scaleMode || "FILL" }];
        parent.appendChild(rect);
        figma.commitUndo();
        return {
          type: request.type,
          requestId: request.requestId,
          data: { id: rect.id, name: rect.name, type: rect.type, bounds: getBounds(rect) }
        };
      }
      default:
        return null;
    }
  });
  const sendStatus = () => {
    figma.ui.postMessage({
      type: "plugin-status",
      payload: {
        fileName: figma.root.name,
        selectionCount: figma.currentPage.selection.length
      }
    });
  };
  const handleRequest = (request) => __async(null, null, function* () {
    var _a;
    try {
      const result = (_a = yield handleReadRequest(request)) != null ? _a : yield handleWriteRequest(request);
      if (result === null)
        throw new Error(`Unknown request type: ${request.type}`);
      return result;
    } catch (error) {
      return {
        type: request.type,
        requestId: request.requestId,
        error: error instanceof Error ? error.message : String(error)
      };
    }
  });
  figma.showUI(__html__, { width: 320, height: 180 });
  sendStatus();
  figma.on("selectionchange", () => {
    sendStatus();
  });
  figma.ui.onmessage = (message) => __async(null, null, function* () {
    if (message.type === "ui-ready") {
      sendStatus();
      return;
    }
    if (message.type === "server-request") {
      const response = yield handleRequest(message.payload);
      try {
        figma.ui.postMessage(response);
      } catch (err) {
        figma.ui.postMessage({
          type: response.type,
          requestId: response.requestId,
          error: err instanceof Error ? err.message : String(err)
        });
      }
    }
  });
})();
