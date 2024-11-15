--@meta
---
--- Run a glue script
---
---@param glue_file string the glue file to run
---
---
function glue(glue_file) end

---
--- Create a runnable group
---
---@param name string the name of the group to run
---@param fn function the function to run when the group is invoked
---
---
function group(name, fn) end

---
--- Creates a backup of a file
---
---@param path string the file to create a backup of
---
---
function backup(path) end

---@class BlockinfileParams
---@field path string the file to insert the block into
---@field block string the multi-line text block to be inserted or updated
---@field insertafter? string the multi-line text block to be inserted or updated
---@field insertbefore? string the multi-line text block to be inserted or updated
---@field marker? string the multi-line text block to be inserted or updated
---@field markerbegin? string the multi-line text block to be inserted or updated
---@field markerend? string the multi-line text block to be inserted or updated
---@field state boolean the multi-line text block to be inserted or updated
---@field backup? boolean the multi-line text block to be inserted or updated
---@field create? boolean the multi-line text block to be inserted or updated


---
--- Insert/update/remove a block of multi-line text surrounded by customizable markers in a file
---
---@param block_params BlockinfileParams the configuration for the block insertion
---
---
function blockinfile(block_params) end

---@class CopyOpts
---@field source string the file or folder to copy
---@field dest string the destination to copy to
---@field strategy? "replace"|"merge" a strategy for how to manage conflicts (defaults to merge)
---@field symlink? "deep"|"shallow"|"skip" how to handle symlinks (copy the content, copy the link, or the default skip)


---
--- Copies folder
---
---@param opts CopyOpts the copy options
---
---
function copy(opts) end

---
--- Installs Homebrew if not already installed
---
---
---
function homebrew_install() end

---
--- Marks a homebrew package for installation
---
---@param pkg string the name of the package to install
---
---
function homebrew(pkg) end

---
--- Marks a cask for installation
---
---@param pkg string the name of the cask to install
---
---
function homebrew_cask(pkg) end

---
--- Marks a homebrew tap for installation
---
---@param tap string the name of the tap to install
---
---
function homebrew_tap(tap) end

---
--- Marks a Mac App Store package for installation
---
---@param name string the name of the mas to install
---
---
function homebrew_mas(name) end

---
--- Marks a whalebrew package for installation
---
---@param name string the name of the whalebrew to install
---
---
function homebrew_whalebrew(name) end

---
--- Installs all marked packages
---
---
---
function homebrew_sync() end

---
--- Upgrades all homebrew packages
---
---
---
function homebrew_upgrade() end

---
--- Run a shell command
---
---@param cmd string the shell command to run
---
---
function sh(cmd) end

---
--- Print a string
---
---@param obj any the message or object to log
---
---
function print(obj) end

---
--- Trims the extra indentation of a multi-line string
---
---@param txt string the text to trim
---
---@return string the trimmed text
---
function trim(txt) end

---
--- Reads a file as a string
---
---@param path string the path of the file to read
---
---@return string the file content
---
function read(path) end


