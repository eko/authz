import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { object, string } from 'yup';

import { AuthenticationExpiration, AuthenticationToken, AuthenticationUser } from 'context/auth';
import { DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue } from 'react-hook-form';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';
import { signIn } from 'service/auth/api';

export type SigninFormData = {
    username: string
    password: string
}

export type OnSubmitHandler = () => FormEventHandler;

type SigninForm = {
    onSubmit: OnSubmitHandler
    register: UseFormRegister<SigninFormData>
    getValues: UseFormGetValues<SigninFormData>
    setValue: UseFormSetValue<SigninFormData>
    defaultValues?: Readonly<DeepPartial<SigninFormData>> | SigninFormData
    errors: FieldErrors<SigninFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<SigninFormData>
}

const schema = object({
    username: string().required('You have to specify a username'),
    password: string().required('You have to specify a password'),
}).required();

export default function useSigninForm(
    navigate: NavigateFunction,
): SigninForm {
    const toast = useToast();

    const {
        register,
        getValues,
        setValue,
        handleSubmit,
        formState: {
          defaultValues,
          errors,
          isSubmitting,
          isValid,
        },
        reset,
        setError,
    } = useForm<SigninFormData>({
        resolver: yupResolver(schema),
        defaultValues: {},
    });

    const onSubmit = () => handleSubmit(async (data: SigninFormData) => {
        const response = await signIn(mapUserFormDataToRequest(data));

        if (isAPIError(response)) {
            setError('username', { type: 'manual' });
            setError('password', { type: 'manual' });

            toast.error(`Unable to authenticate: ${response.message}`);
            return;
        }

        const expireAt = new Date();
        expireAt.setSeconds(expireAt.getSeconds() + response.expires_in);

        localStorage.setItem(AuthenticationExpiration, expireAt.toISOString());
        localStorage.setItem(AuthenticationToken, response.access_token);
        localStorage.setItem(AuthenticationUser, JSON.stringify(response.user));

        toast.success(`You're now authenticated.`);

        navigate('/');
    });

    return {
        onSubmit,
        register,
        getValues,
        setValue,
        defaultValues,
        isSubmitting,
        isValid,
        errors,
        reset,
    };
}

export const mapUserFormDataToRequest = (data: SigninFormData): any => {
    return {
        username: data.username,
        password: data.password,
    };
}