test("test 1", function()
    print("Test 1 is running")
    assert(true, "value should have been true")
    return nil
end)

test("test 2", function()
    print("Test 2 is running")
    assert(false, "value should have been true")
    return nil
end)
