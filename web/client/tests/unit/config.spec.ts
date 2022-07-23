// import { shallowMount } from "@vue/test-utils";
// import HelloWorld from "@/components/HelloWorld.vue";
import config from "../../src/config";

describe("HelloWorld.vue", () => {
  it("renders props.msg when passed", () => {
    // const msg = "new message";
    // const wrapper = shallowMount(HelloWorld, {
    //   propsData: { msg }
    // });
    expect(config.appName).toBe("PhoneInfoga");
    expect(config.appDescription).toBe(
      "Advanced information gathering & OSINT tool for phone numbers"
    );
    expect(config.apiUrl).toBe("/api");
  });
});
