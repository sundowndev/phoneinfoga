import { formatNumber, isValid } from "../../src/utils";

describe("src/utils", () => {
  describe("#formatNumber", () => {
    it("should format phone number", () => {
      expect(formatNumber("+254 (74_370-6303")).toBe("254743706303");
      expect(formatNumber("1/(254)_7437-06303")).toBe("254743706303");
    });
  });

  describe("#isValid", () => {
    it("should check if number is valid", () => {
      expect(isValid("+254/(743)_706-303")).toBe(true);
      expect(isValid("254 743 706303")).toBe(true);
      expect(isValid("this254 743 706303")).toBe(false);
    });
  });
});
