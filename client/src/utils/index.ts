const formatNumber = (number: string): string => {
  return number.replace(/[_\W]+/g, "");
};

const isValid = (number: string): boolean => {
  const formatted = formatNumber(number);

  return formatted.match(/^[0-9]+$/) !== null && formatted.length > 2;
};

export { formatNumber, isValid };
