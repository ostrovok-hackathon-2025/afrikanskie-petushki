"use client";

import { useCallback, useMemo, useState } from "react"
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { getSecretGuestAPI } from "@/api/api";
import { signIn } from "next-auth/react";
import { redirect } from "next/navigation";
import { ERR_EMPTY, ERR_INVALID_VALUE, ERR_LATIN_ONLY, ERR_NO_CAPITALS, ERR_NO_DIGITS, ERR_NO_SPECIAL_CHARACTERS, ERR_TOO_LONG, ERR_TOO_SHORT, NO_ERR, useForm } from "@/lib/formHook/form-hook";
import { cn } from "@/lib/utils";

const { postUserLogIn } = getSecretGuestAPI();

export default function SignUp() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const [validate, getError, markInvalid, withReset] = useForm({
        username: {
            type: "text",
            notEmpty: true,
            minLength: 3,
            maxLength: 32,
            latinOnly: true
        },
        
        password: {
            type: "password",
            includeDigits: true,
            includeCapitals: true,
            includeSpecialCharacters: true,
            minLength: 8,
            maxLength: 64,
        },
        
    }); 

    const handleLogIn = useCallback(async () => {
        if (
            !validate({
                username: username,
                password: password
            })
        ) {
            return;
        }

        try {
            const authData = await postUserLogIn({
                ostrovok_login: username,
                password: password
            });

            await signIn("credentials", {
                redirect: false,
                id: "",
                accessToken: authData.data.access_token ?? "",
                refreshToken: authData.data.refresh_token ?? "",
                accessTTL: authData.data.access_ttl?.toString() ?? "",
                refresh_ttl: authData.data.refresh_ttl?.toString() ?? "",
            });

            redirect("/");
        } catch (err) {
            console.error(err);
        }
    }, [username, password]);

    const errorText = useMemo(() => {
        const usernameErr = getError("username");

        switch (usernameErr) {
        case ERR_EMPTY:
            return "введите имя";
        case ERR_TOO_SHORT:
            return "слишком короткое имя";
        case ERR_TOO_LONG:
            return "слишком длинное имя";
        case ERR_LATIN_ONLY:
            return "имя может содержать только буквы латинского алфавита, цифры и _";
        case ERR_INVALID_VALUE:
            return "пользователь не найден";
        }

        const passwordErr = getError("password");

        switch (passwordErr) {
        case ERR_EMPTY:
            return "введите пароль";
        case ERR_TOO_SHORT:
            return "слишком короткий пароль";
        case ERR_TOO_LONG:
            return "слишком длинный пароль";
        case ERR_NO_CAPITALS:
            return "пароль должен содержать заглавные буквы";
        case ERR_NO_DIGITS:
            return "пароль должен содержать цифры";
        case ERR_LATIN_ONLY:
            return "пароль может содержать только буквы латинского алфавита, цифры и спец. символы";
        case ERR_NO_SPECIAL_CHARACTERS:
            return "пароль должен содержать спец. символ";
        case ERR_INVALID_VALUE:
            return "неправильный пароль";
        }

        return "no";
    }, [getError("username"), getError("password")]);

    return <div className="w-full h-full flex items-center justify-center">
        <div className="w-1/3 h-1/2 rounded-lg border border-foreground box-border p-5">
            <h2 className="font-gain text-xl mb-6">Вход</h2>

            <div className="font-gain text-base mb-1">Логин Островка</div>
            <Input 
                value={username} 
                onChange={withReset((e) => setUsername(e.target.value))} 
                placeholder="логин" 
                className={cn("mb-2", getError("username") !== NO_ERR && "border-destructive")} 
                type="text"
            />
            
            <div className="font-gain text-base mb-1">Пароль</div>
            <Input 
                value={password} 
                onChange={withReset((e) => setPassword(e.target.value))} 
                placeholder="password" 
                className={cn("mb-3", getError("password") !== NO_ERR && "border-destructive")}  
                type="password"
            />

            <div className="mb-4">
                Нет аккаунта? <a href="/sign-up" className="text-primary">Зарегистрироваться</a>
            </div>

            <div
                className={cn(
                "-full text-xs font-unbounded text-destructive mb-4",
                errorText === "no" && "opacity-0"
                )}
            >
                {errorText}
            </div>
            
            <Button onClick={handleLogIn} size={"default"} className="mb-2">Войти</Button>

        </div>
    </div>
}