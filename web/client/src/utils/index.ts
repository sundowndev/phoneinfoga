import axios from "axios";
import config from "../config/index";
interface ScannerObject {
  name: string;
  description: string;
}

const formatNumber = (number: string): string => {
  return number.replace(/[_\W]+/g, "");
};

const isValid = (number: string): boolean => {
  const formatted = formatNumber(number);

  return formatted.match(/^[0-9]+$/) !== null && formatted.length > 2;
};

const formatString = (string: string): string => {
  return string.replace(/([A-Z])/g, " $1").trim();
};

const getScanners = async (): Promise<ScannerObject[]> => {
  const res = await axios.get(`${config.apiUrl}/v2/scanners`);

  // TODO: Remove this filter once the scanner local is remove
  return res.data.scanners.filter(
    (scanner: ScannerObject) => scanner.name !== "local"
  );
};

export { formatNumber, isValid, formatString, getScanners };
