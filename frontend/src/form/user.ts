import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { object, string } from 'yup';

import { Control, DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue } from 'react-hook-form';
import { createUser } from 'service/model/user';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';
import { User } from 'service/model/model';

export type UserFormData = {
    username: string
}

export type OnSubmitHandler = (token: string, callback: (user: User) => void) => FormEventHandler;

type UserForm = {
    control: Control<UserFormData>
    onSubmit: OnSubmitHandler
    register: UseFormRegister<UserFormData>
    getValues: UseFormGetValues<UserFormData>
    setValue: UseFormSetValue<UserFormData>
    defaultValues?: Readonly<DeepPartial<UserFormData>> | UserFormData
    errors: FieldErrors<UserFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<UserFormData>
}

const schema = object({
    username: string().required('You have to specify a username.'),
}).required();

export default function useUserForm(
    navigate: NavigateFunction,
    user?: User,
): UserForm {
    const toast = useToast();

    const {
        control,
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
    } = useForm<UserFormData>({
        resolver: yupResolver(schema),
        defaultValues: mapUserToFormData(user),
    });

    const onSubmit = (token: string, callback: (user: User) => void) => handleSubmit(async (data: UserFormData) => {
        if (user !== undefined) {
            return;
        }

        const response = await createUser(token, mapUserFormDataToRequest(data));

        if (isAPIError(response)) {
            toast.error(`Unable to create user: ${response.message}`);
        } else {
            toast.success(`User ${response.username} has been created.`);
            callback(response);
        }
    });

    return {
        control,
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

export const mapUserToFormData = (data?: User): UserFormData => {
    return ({
        username: data?.username || '',
    });
}

export const mapUserFormDataToRequest = (data: UserFormData): any => {
    return {
        username: data.username,
    };
}