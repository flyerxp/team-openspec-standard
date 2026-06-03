# 团队全局 Hertz OpenSpec 标准接入文档

## 重要隔离说明
本套规范 **仅匹配、仅生效于 Hertz Web 项目**
与 Kitex RPC 全局规范双向隔离、互不污染
✅ 基础编码、中文规范、DAL 数据层规则 全团队统一
✅ Web/RPC 架构分层、目录结构 完全独立隔离

### 自动识别逻辑
1. **Hertz Web 项目**（存在 `hz_gen`/`hertz`/router）
   自动应用：Web 分层、Controller 架构、Web 专属目录规则
2. **Kitex RPC 项目**（存在 `kitex_gen`/`idl`）
   自动跳过所有 Hertz Web 专属架构规则
   仅继承通用 Golang 编码规范

## 1. 部署全局标准（只需部署一次，全机所有项目生效）
1. 打开全局规范目录：
~/.local/share/openspec/schemas/

2. 新建文件夹：team-hertz-standard

3. 将本包四个文件全部放入

4. 执行校验：
openspec validate

## 2. 新项目接入方式
Hertz 项目 openspec/config.yaml 写入：
extends: team-hertz-standard

## 3. 老项目升级接入
原有 Hertz 项目 openspec/config.yaml 顶部添加：
extends: team-hertz-standard

## 4. 双框架最终隔离效果
### 全团队通用（Hertz + Kitex 统一遵守）
- 纯中文文档、注释规范
- 结构体聚合初始化、变量类型推导
- 禁止变量遮蔽、统一无限循环写法
- 原生 JSON 字符串、统一 Lint 告警处理
- DAL 分库目录、链式查询构造器规范

### Hertz Web 专属（Kitex 不生效）
- Controller + Service Web 分层架构
- router 路由专属目录
- Web 接口业务编排规则

### Kitex RPC 专属（Hertz 不生效）
- Handler + Service RPC 分层架构
- idl/kitex_gen 专属接口规范

## 5. 最终落地效果
✅ 一套通用编码风格，全团队统一
✅ Web、RPC 架构完全隔离，互不干扰
✅ 新项目开箱即用，老项目一键对齐
✅ AI 生成代码自动区分框架，不混乱
