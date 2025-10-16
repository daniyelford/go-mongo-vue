<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\File;

class MakeHttpStub extends Command
{
    protected $signature = 'make:http-stub {--force : بازسازی در صورت وجود}';
    protected $description = 'ایجاد ساختار پایه Http Request و Response برای ارتباط با Go API';

    public function handle()
    {
        $basePaths = [
            app_path('Http/Request'),
            app_path('Http/Response'),
        ];

        foreach ($basePaths as $path) {
            if (!File::exists($path)) {
                File::makeDirectory($path, 0755, true);
                $this->info("📁 پوشه ساخته شد: $path");
            }
        }

        // ایجاد VueRequest
        $vueRequestPath = app_path('Http/Request/VueRequest.php');
        if (!File::exists($vueRequestPath) || $this->option('force')) {
            File::put($vueRequestPath, <<<PHP
<?php

namespace App\Http\Request;

use Illuminate\Support\Facades\DB;

class VueRequest
{
    public static function send(array \$data)
    {
        DB::table('api_requests')->insert([
            'action' => \$data['action'] ?? null,
            'payload' => json_encode(\$data['payload'] ?? []),
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        // در اینجا می‌توانی درخواست HTTP را به Go ارسال کنی
    }
}
PHP);
            $this->info("✅ فایل ساخته شد: $vueRequestPath");
        }

        // ایجاد GoResponse
        $goResponsePath = app_path('Http/Response/GoResponse.php');
        if (!File::exists($goResponsePath) || $this->option('force')) {
            File::put($goResponsePath, <<<PHP
<?php

namespace App\Http\Response;

use Illuminate\Support\Facades\DB;

class GoResponse
{
    public static function handle(\$response)
    {
        DB::table('api_responses')->insert([
            'response' => json_encode(\$response),
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        return \$response;
    }
}
PHP);
            $this->info("✅ فایل ساخته شد: $goResponsePath");
        }

        $this->info('🎉 ساختار کامل Http Request و Response آماده است!');
    }
}
