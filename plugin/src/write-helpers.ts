// Write helpers — utilities used exclusively by write handlers.

export const hexToRgb = (hex: string) => {
  const clean = hex.replace("#", "");
  return {
    r: parseInt(clean.slice(0, 2), 16) / 255,
    g: parseInt(clean.slice(2, 4), 16) / 255,
    b: parseInt(clean.slice(4, 6), 16) / 255,
    a: clean.length >= 8 ? parseInt(clean.slice(6, 8), 16) / 255 : 1,
  };
};

export const makeSolidPaint = (colorInput: any, opacityOverride?: number): SolidPaint => {
  const { r, g, b, a } = typeof colorInput === "string"
    ? hexToRgb(colorInput)
    : { r: colorInput.r, g: colorInput.g, b: colorInput.b, a: colorInput.a != null ? colorInput.a : 1 };
  const eff = opacityOverride != null ? opacityOverride : a;
  const paint: any = { type: "SOLID", color: { r, g, b } };
  if (eff !== 1) paint.opacity = eff;
  return paint;
};

export const getParentNode = async (parentId: string | undefined) => {
  if (!parentId) return figma.currentPage;
  const parent = await figma.getNodeByIdAsync(parentId);
  if (!parent) throw new Error(`Parent node not found: ${parentId}`);
  if (!("appendChild" in parent)) throw new Error(`Node ${parentId} cannot have children`);
  return parent as ChildrenMixin & BaseNode;
};

export const base64ToBytes = (b64: string) => {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
  const lookup: Record<string, number> = {};
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
    bytes[j++] = (a << 2) | (bv >> 4);
    if (j < outLen) bytes[j++] = ((bv & 15) << 4) | (c >> 2);
    if (j < outLen) bytes[j++] = ((c & 3) << 6) | d;
  }
  return bytes;
};

export const makeRGBA = (colorInput: any, opacityOverride?: number): RGBA => {
  const { r, g, b, a } = typeof colorInput === "string"
    ? hexToRgb(colorInput)
    : { r: colorInput.r, g: colorInput.g, b: colorInput.b, a: colorInput.a != null ? colorInput.a : 1 };
  return { r, g, b, a: opacityOverride != null ? opacityOverride : a };
};

export const loadAllFonts = async (node: TextNode) => {
  if (node.hasMissingFont) {
    throw new Error(`Text node ${node.id} has missing fonts`);
  }

  if (typeof node.fontName !== "symbol" && node.fontName) {
    await figma.loadFontAsync(node.fontName);
    return;
  }

  const seen = new Set<string>();
  const fonts = node.getRangeAllFontNames(0, node.characters.length);
  for (const font of fonts) {
    const key = `${font.family}::${font.style}`;
    if (seen.has(key)) continue;
    seen.add(key);
    await figma.loadFontAsync(font);
  }
};

export const parseLineHeight = (value: any): LineHeight | undefined => {
  if (!value || typeof value !== "object") return undefined;
  if (value.unit === "AUTO") return { unit: "AUTO" };
  if (!value.unit || value.value == null) return undefined;
  return { unit: value.unit, value: value.value };
};

export const parseLetterSpacing = (value: any): LetterSpacing | undefined => {
  if (!value || typeof value !== "object") return undefined;
  if (!value.unit || value.value == null) return undefined;
  return { unit: value.unit, value: value.value };
};

export const applyTextBoxProperties = (
  node: TextNode,
  value: { textAutoResize?: TextAutoResize; width?: number; height?: number },
  applied?: Record<string, any>,
) => {
  const hasWidth = value.width != null;
  const hasHeight = value.height != null;

  let nextMode = value.textAutoResize;
  if (!nextMode && (hasWidth || hasHeight)) {
    nextMode = hasWidth && !hasHeight ? "HEIGHT" : "NONE";
  }

  if (nextMode) {
    node.textAutoResize = nextMode;
    if (applied) applied.textAutoResize = nextMode;
  }

  if (!hasWidth && !hasHeight) return;

  const nextWidth = hasWidth ? value.width! : node.width;
  const nextHeight = hasHeight ? value.height! : node.height;
  node.resize(nextWidth, nextHeight);

  if (applied && hasWidth) applied.width = nextWidth;
  if (applied && hasHeight) applied.height = nextHeight;
};

export const makeEffect = (effect: any): Effect => {
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
        showShadowBehindNode: effect.showShadowBehindNode === true,
      };
    case "INNER_SHADOW":
      return {
        type: "INNER_SHADOW",
        color: makeRGBA(effect.color),
        offset: { x: effect.x != null ? effect.x : 0, y: effect.y != null ? effect.y : 2 },
        radius: effect.radius != null ? effect.radius : 12,
        spread: effect.spread != null ? effect.spread : 0,
        visible: effect.visible != null ? effect.visible : true,
        blendMode: effect.blendMode || "NORMAL",
      };
    case "LAYER_BLUR":
    case "BACKGROUND_BLUR":
      return {
        type: effect.type,
        radius: effect.radius != null ? effect.radius : 12,
        visible: effect.visible != null ? effect.visible : true,
        blurType: "NORMAL",
      } as Effect;
    default:
      throw new Error(`Unsupported effect type: ${effect.type}`);
  }
};

const toComponentNode = (node: BaseNode | null): ComponentNode | null => {
  if (!node) return null;
  if (node.type === "COMPONENT") return node;
  if (node.type === "COMPONENT_SET") return node.defaultVariant;
  return null;
};

export const resolveComponentNode = async (selector: {
  componentId?: string;
  componentKey?: string;
  name?: string;
}): Promise<ComponentNode> => {
  if (selector.componentId) {
    const node = await figma.getNodeByIdAsync(selector.componentId);
    const component = toComponentNode(node);
    if (component) return component;
    throw new Error(`Component not found: ${selector.componentId}`);
  }

  const pages = figma.root.children;
  let caseInsensitiveMatch: ComponentNode | null = null;

  for (const page of pages) {
    await page.loadAsync();
    const candidates = page.findAllWithCriteria({
      types: ["COMPONENT", "COMPONENT_SET"],
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
};

export const mapComponentProperties = (
  instance: InstanceNode,
  rawProps: Record<string, any>,
): { [propertyName: string]: string | boolean } => {
  const mapped: { [propertyName: string]: string | boolean } = {};
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
