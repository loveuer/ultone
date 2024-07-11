# utl-one: utl all in one

### Usage

#### 第一步: 克隆项目

- 1. `git clone -b master --depth 1 https://github.com/loveuer/ultone.git {your_project_name}`
  
- 2. `cd {your_project_name} && rm -rf .git && git init`


#### 第二步: 设置为你自己的项目名称

- 1. `go mod edit -module {your_project_name}`

   - 2. `修改 go 文件中的 module 名称`
      * `Windows 下使用 powershell`
        ```psh
        $NEW_MODULE_NAME = "{your_project_name}"
        $OLD_MODULE = "ultone"

        Get-ChildItem -Path . -Filter '*.go' -Recurse | ForEach-Object {
            $content = Get-Content -Path $_.FullName
            $updatedContent = $content -replace $OLD_MODULE, $NEW_MODULE_NAME
            $updatedContent | Set-Content -Path $filePath
        }
        ```
      
      * `MacOS 下`
        ```sh
        find . -type f -name '*.go' -exec sed -i '' -e 's/ultone/{your_project_name}/g' {} \;
        ```
      * `Linux 下`  
        ```sh
        find . -type f -name '*.go' -exec sed -i -e 's,ultone,{your_project_name},g' {} \;
        ```
   
- 3. `go mod tidy`

- 4. 初始账号: `admin`, 初始密码: `123456`

### Setting

#### 仔细查看项目中的 todo 

#### 仔细查看 opt.var 中的设置

#### SQL

- sqlite:
- postgresql:
- mysql

#### Cache

- redis
- memory

### Feature

- 用户全功能模块
- 操作日志

### Next

- common user list (比如操作日志用户下拉)