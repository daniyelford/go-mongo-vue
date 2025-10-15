<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\File;

class MakeHttpStub extends Command
{
    protected $signature = 'make:http-stub';
    protected $description = 'Create Http/Request and Http/Response folders with base classes';

    public function handle()
    {
        $folders = [
            app_path('Http/Request'),
            app_path('Http/Response'),
        ];

        foreach ($folders as $folder) {
            if (!File::exists($folder)) {
                File::makeDirectory($folder, 0755, true);
                $this->info("Folder created: $folder");
            } else {
                $this->info("Folder already exists: $folder");
            }
        }

        $requestFile = app_path('Http/Request/VueRequest.php');
        if (!File::exists($requestFile)) {
            File::put($requestFile, "<?php\n\nnamespace App\Http\Request;\n\nclass VueRequest\n{\n    public static function send(array \$data)\n    {\n        // Here you can send data to Go server using HTTP client\n    }\n}");
            $this->info("File created: $requestFile");
        }

        $responseFile = app_path('Http/Response/GoResponse.php');
        if (!File::exists($responseFile)) {
            File::put($responseFile, "<?php\n\nnamespace App\Http\Response;\n\nclass GoResponse\n{\n    public static function handle(\$response)\n    {\n        // Here you can format response from Go\n        return \$response;\n    }\n}");
            $this->info("File created: $responseFile");
        }

        $this->info('Http stub generated successfully!');
    }
}
