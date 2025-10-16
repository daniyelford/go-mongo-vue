<?php

namespace App\Http\Requests\Api\Auth;

use Illuminate\Foundation\Http\FormRequest;

class SendCodeRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true; // اگر نیاز به احراز خاص نداری
    }

    public function rules(): array
    {
        return [
            'mobile'  => 'required|regex:/^[0-9]{10,15}$/',
            'country' => 'required|string|size:2', // مثل "98"
        ];
    }

    public function messages(): array
    {
        return [
            'mobile.required' => 'شماره موبایل الزامی است.',
            'mobile.regex'    => 'فرمت شماره موبایل نادرست است.',
            'country.required' => 'کد کشور الزامی است.',
        ];
    }
}
