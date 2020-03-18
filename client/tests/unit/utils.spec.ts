import { formatNumber, isValid } from "../../src/utils";

describe("src/utils", () => {
  describe("#formatNumber", () => {
    it("should format phone number", () => {
      expect(formatNumber("+1 (555) 444-3333")).toBe("15554443333");
      expect(formatNumber("1/(555)_444-3333")).toBe("15554443333");
    });
  });

  describe("#isValid", () => {
    it("should check if number is valid", () => {
      expect(isValid("+1/(555)_444-3333")).toBe(true);
      expect(isValid("1 555 4443333")).toBe(true);
      expect(isValid("this1 555 4443333")).toBe(false);
    });
  });
});
