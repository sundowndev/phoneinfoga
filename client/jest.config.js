module.exports = {
    // moduleFileExtensions: [
    //   'js',
    //   'jsx',
    //   'json',
    //   'vue',
    //   'ts',
    //   'tsx'
    // ],
    // "preset": "@vue/cli-plugin-unit-jest/presets/typescript",
    transform: {
      '^.+\\.vue$': 'vue-jest',
      '.+\\.(css|styl|less|sass|scss|svg|png|jpg|ttf|woff|woff2)$': 'jest-transform-stub',
      '^.+\\.tsx?$': 'ts-jest'
    },
    // transformIgnorePatterns: [
    //   '/node_modules/'
    // ],
    // moduleNameMapper: {
    //   '^@/(.*)$': '<rootDir>/src/$1'
    // },
    // snapshotSerializers: [
    //   'jest-serializer-vue'
    // ],
    testMatch: [
      '**/tests/unit/**/*.spec.(js|jsx|ts|tsx)|**/__tests__/*.(js|jsx|ts|tsx)'
    ],
    // testURL: 'http://localhost/',
    // watchPlugins: [
    //   'jest-watch-typeahead/filename',
    //   'jest-watch-typeahead/testname'
    // ],
    // globals: {
    //   'ts-jest': {
    //     babelConfig: true
    //   }
    // },
    collectCoverage: true,
    "collectCoverageFrom": [
      "src/(store|config|components|views|utils|models)/**/*.(ts|vue)"
    ]
  }  