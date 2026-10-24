// Wails 绑定会把 Go 端 error 透传到前端，可能是字符串也可能是 Error。
// 401 文案历史上有两种："401 Unauthorized"（旧版 errors.New）和
// "HTTP 401 from <url>"（HTTPStatusError），两种都要能识别。
export function isUnauthorized(error: unknown): boolean {
    const msg = typeof error === "string" ? error : error instanceof Error ? error.message : String(error);
    return msg.startsWith("HTTP 401") || msg.includes("401 Unauthorized");
}
