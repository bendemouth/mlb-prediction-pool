function formatter(value: number, decimals: number = 2): string {
  return value.toFixed(decimals);
}

export default formatter;