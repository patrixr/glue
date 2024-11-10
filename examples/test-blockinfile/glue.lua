blockinfile({
    path = "./sample.txt",
    state = true,
    insertafter = "line 2",
    block = trim([[
        This is the content to insert
    ]])
})
