#[test]
fn duplicate_resource_name_attribute_fails() {
    let t = trybuild::TestCases::new();
    t.compile_fail("tests/ui/duplicate_resource_name_attr.rs");
}
