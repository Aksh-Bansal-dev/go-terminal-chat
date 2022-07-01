const colorMap = new Map<number, string>();

export const getColor = (colorCode: number): string => {
  if (colorMap.has(colorCode)) return colorMap.get(colorCode)!;
  const randColor = Math.floor(Math.random() * 16777215).toString(16);
  colorMap.set(colorCode, randColor);
  return randColor;
};
