<?php

namespace App\Helpers;

use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use Illuminate\Support\Facades\Redis;

class JwtHelper
{
    /**
     * ساخت JWT Token و ذخیره در Redis
     */
    public static function createToken(string $mobile): string
    {
        $payload = [
            'iss' => config('app.url'),    // صادرکننده
            'iat' => time(),               // زمان ایجاد
            'exp' => time() + 3600,        // زمان انقضا (۱ ساعت)
            'mobile' => $mobile,           // اطلاعات کاربر
        ];

        $secret = env('JWT_SECRET');

        // تولید توکن
        $token = JWT::encode($payload, $secret, 'HS256');

        // ذخیره در Redis با زمان انقضا ۱ ساعت
        Redis::setex("token:$mobile", 3600, $token);

        return $token;
    }

    /**
     * بررسی و decode توکن (برای middleware)
     */
    public static function verifyToken(string $token): ?object
    {
        try {
            $decoded = JWT::decode($token, new Key(env('JWT_SECRET'), 'HS256'));
            return $decoded;
        } catch (\Exception $e) {
            return null;
        }
    }
}
