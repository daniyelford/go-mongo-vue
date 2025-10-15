<?php

namespace App\Http\Response;

use Illuminate\Support\Facades\DB;

class GoResponse
{
    public static function handle($response)
    {
        // ذخیره پاسخ در جدول api_responses
        DB::table('api_responses')->insert([
            'response' => json_encode($response),
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        return $response;
    }
}
