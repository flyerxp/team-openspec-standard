# 团队全局 Kitex OpenSpec 标准接入文档

## 重要说明（解决混用问题）
本套规范 **仅匹配、仅生效于 Kitex RPC 项目**
**完全不会污染 Hertz Web 项目**，双框架隔离生效

### 自动识别逻辑
1. **Kitex 项目**（存在 `kitex_gen`/`idl`/`kitex.yaml`）
   自动应用：RPC 分层、DAL分库目录、Service架构、RPC专属规则
2. **Hertz Web 项目**（存在 `hz_gen`/`hertz` 标识）
   **自动跳过所有 Kitex 专属架构规则**
   仅继承通用 Golang 编码规范（结构体初始化、变量推导、字符串规范等）

## 1. 部署全局标准（只需部署一次，全机所有项目生效）
1. 打开目录：
~/.local/share/openspec/schemas/

2. 新建文件夹：team-kitex-standard

3. 将本包四个文件全部放入

4. 执行校验：
openspec validate

## 2. 新项目接入方式
项目根目录 openspec/config.yaml 写入：
extends: team-kitex-standard

- Kitex项目：自动完整继承RPC架构+编码规范
- Hertz项目：**只继承通用编码规范，不继承RPC分层/目录规则**

## 3. 老项目升级接入
直接在原有 openspec/config.yaml 顶部添加：
extends: team-kitex-standard

## 4. 双框架隔离生效范围
### 全局通用（Kitex + Hertz 全部遵守）
- 中文注释/中文文档
- 结构体聚合初始化
- 变量类型推导、禁止变量遮蔽
- 原生JSON字符串、Lint注释规范

### 仅 Kitex RPC 专属（Hertz 不生效）
- Handler/Service/Logic/DAL/Convert 分层架构
- 分库DAL目录结构 `dal/gormL/{db_name}/where`
- RPC Service 结构体模板
- RPC专属依赖校验规则

## 5. 最终效果
✅ Kitex RPC项目：全套专属RPC架构规范
✅ Hertz Web项目：只统一编码风格，**不套用RPC目录/分层**
✅ 双框架共存、互不干扰、全局统一编码风格
✅ 彻底解决规范互相污染问题
