module.exports = {
  transform: {
    "^.+\\.vue$": "vue-jest",
    ".+\\.(css|styl|less|sass|scss|svg|png|jpg|ttf|woff|woff2)$":
      "jest-transform-stub",
    "^.+\\.tsx?$": "ts-jest",
  },
  testMatch: [
    "**/tests/unit/**/*.spec.(js|jsx|ts|tsx)|**/__tests__/*.(js|jsx|ts|tsx)",
  ],
  collectCoverage: true,
  collectCoverageFrom: [
    "src/(store|config|components|views|utils|models)/**/*.(ts|vue)",
  ],
};
