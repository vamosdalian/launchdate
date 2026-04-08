export default function MobileWarningPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background px-6 text-center text-foreground">
      <div className="max-w-sm space-y-4">
        <h1 className="text-2xl font-semibold">请使用桌面端访问</h1>
        <p className="text-sm text-muted-foreground">
          当前管理后台仅支持在电脑端使用。请使用电脑浏览器访问 admin.launch-date.com 以继续操作。
        </p>
      </div>
    </div>
  );
}
