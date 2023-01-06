import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { object, string } from 'yup';

import { Control, DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue } from 'react-hook-form';
import { createClient } from 'service/model/client';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';
import { Client } from 'service/model/model';

export type ClientFormData = {
    name: string
}

export type OnSubmitHandler = (token: string) => FormEventHandler;

type ClientForm = {
    control: Control<ClientFormData>
    onSubmit: OnSubmitHandler
    register: UseFormRegister<ClientFormData>
    getValues: UseFormGetValues<ClientFormData>
    setValue: UseFormSetValue<ClientFormData>
    defaultValues?: Readonly<DeepPartial<ClientFormData>> | ClientFormData
    errors: FieldErrors<ClientFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<ClientFormData>
}

const schema = object({
    name: string().required('You have to specify a name.'),
}).required();

export default function useClientForm(
    navigate: NavigateFunction,
    client?: Client,
): ClientForm {
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
    } = useForm<ClientFormData>({
        resolver: yupResolver(schema),
        defaultValues: mapClientToFormData(client),
    });

    const onSubmit = (token: string) => handleSubmit(async (data: ClientFormData) => {
        if (client !== undefined) {
            return;
        }

        const response = await createClient(token, mapClientFormDataToRequest(data));

        if (isAPIError(response)) {
            toast.error(`Unable to create service account: ${response.message}`);
        } else {
            toast.success(`Service account ${response.name} has been created.`);
            navigate(`/clients/edit/${response.client_id}`);
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

export const mapClientToFormData = (data?: Client): ClientFormData => {
    return ({
        name: data?.name || '',
    });
}

export const mapClientFormDataToRequest = (data: ClientFormData): any => {
    return {
        name: data.name,
    };
}