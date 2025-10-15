<?php

namespace App\Http\Request;

use Illuminate\Support\Facades\DB;

class VueRequest
{
    public static function send(array $data)
    {
        // ذخیره داده‌ها در جدول api_requests
        DB::table('api_requests')->insert([
            'action' => $data['action'] ?? null,
            'payload' => json_encode($data['payload'] ?? []),
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        // اینجا می‌تونی کد ارسال HTTP به سرور Go هم اضافه کنی
    }
}
