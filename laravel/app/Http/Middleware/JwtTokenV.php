<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use App\Helpers\JwtHelper;
use Illuminate\Support\Facades\Redis;


class JwtTokenV
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next)
    {
        $authHeader = $request->header('Authorization');

        // بررسی header Authorization
        if (!$authHeader || !str_starts_with($authHeader, 'Bearer ')) {
            return response()->json(['message' => 'missing token'], 401);
        }

        $tokenString = substr($authHeader, 7);

        try {
            // بررسی توکن با JwtHelper
            $credentials = JwtHelper::verifyToken($tokenString);

            if (!$credentials || empty($credentials->mobile)) {
                return response()->json(['message' => 'invalid token claims'], 401);
            }

            $mobile = $credentials->mobile;

            // بررسی Redis برای اطمینان از معتبر بودن token
            $val = Redis::get("token:$mobile");
            if (!$val || $val !== $tokenString) {
                return response()->json(['message' => 'expired or revoked token'], 401);
            }

            // ست کردن موبایل در request برای کنترلرها
            $request->attributes->set('mobile', $mobile);

        } catch (\Exception $e) {
            return response()->json([
                'message' => 'invalid token',
                'error' => $e->getMessage()
            ], 401);
        }

        return $next($request);
    }
}
