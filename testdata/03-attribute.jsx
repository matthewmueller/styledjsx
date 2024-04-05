export default () => {
  const Element = "div"
  return (
    <div>
      <div class="test" {...test.test} />
      <div class="test" {...test.test.test} />
      <div class="test" {...this.test.test} />
      <div data-test="test" />
      <div class="test" />
      <div class={"test"} />
      <div class={`test`} />
      <div class={`test${true ? " test2" : ""}`} />
      <div class={"test " + test} />
      <div class={["test", "test2"].join(" ")} />
      <div class={true && "test"} />
      <div class={test ? "test" : null} />
      <div class={test} />
      <div class={test && "test"} />
      <div class={test && test("test")} />
      <div class={undefined} />
      <div class={null} />
      <div class={false} />
      <div class={"test"} data-test />
      <div data-test class={"test"} />
      <div class={"test"} data-test="test" />
      <div class={"test"} {...props} />
      <div class={"test"} {...props} {...rest} />
      <div class={`test ${test ? "test" : ""}`} {...props} />
      <div class={test && test("test")} {...props} />
      <div class={test && test("test") && "test"} {...props} />
      <div class={test && test("test") && test2("test")} {...props} />
      <div {...props} class={"test"} />
      <div {...props} {...rest} class={"test"} />
      <div {...props} class={"test"} {...rest} />
      <div {...props} />
      <div {...props} {...rest} />
      <div {...props} data-foo {...rest} />
      <div {...props} class={"test"} data-foo {...rest} />
      <div {...{ id: "foo" }} />
      <div {...{ class: "foo" }} />
      <div {...{ class: "foo" }} class="test" />
      <div class="test" {...{ class: "foo" }} />
      <div {...{ class: "foo" }} {...bar} />
      <div {...{ class: "foo" }} {...bar} class="test" />
      <div class="test" {...{ class: "foo" }} {...bar} />
      <div class="test" {...{ class: props.class }} />
      <div class="test" {...{ class: props.class }} {...bar} />
      <div class="test" {...bar} {...{ class: props.class }} />
      <div class="test" {...bar()} />
      <Element />
      <Element class="test" />
      <Element {...props} />
      <style scoped>{"div { color: red }"}</style>
    </div>
  )
}
