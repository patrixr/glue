print("root level")

group("group1", function()
    print("group1 is running")
end)

group("group2", function()
    print("group2 is running")

    group("internal", function()
        print("group2.internal is running")
    end)
end)
