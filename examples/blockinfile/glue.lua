blockinfile({
    state = true,
    insertafter = "line 2",
    path = "sample.txt",
    block = "This is the content to insert",
    backup = true
})
