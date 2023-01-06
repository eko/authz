import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { mixed, object } from 'yup';

import { Control, DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { ItemType } from 'component/MultipleAutocompleteInput';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';
import { check, CheckResponse } from 'service/model/check';

export type CheckFormData = {
    action: ItemType
    principal: ItemType
    resource: ItemType
}

export type OnSubmitHandler = (token: string, callback: (data: CheckResponse) => void) => FormEventHandler;

type CheckForm = {
    control: Control<CheckFormData>
    onSubmit: OnSubmitHandler
    register: UseFormRegister<CheckFormData>
    getValues: UseFormGetValues<CheckFormData>
    setValue: UseFormSetValue<CheckFormData>
    defaultValues?: Readonly<DeepPartial<CheckFormData>> | CheckFormData
    errors: FieldErrors<CheckFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<CheckFormData>
    watch: UseFormWatch<CheckFormData>
}

const schema = object({
    action: mixed<ItemType>().required('This field is required.'),
    resource: mixed<ItemType>().required('This field is required.'),
    principal: mixed<ItemType>().required('This field is required.'),
}).required();

export default function useCheckForm(
    navigate: NavigateFunction,
): CheckForm {
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
        watch,
    } = useForm<CheckFormData>({
        resolver: yupResolver(schema),
    });

    const onSubmit = (token: string, callback: (data: CheckResponse) => void) => handleSubmit(async (data: CheckFormData) => {
        const response = await check(token, mapCheckFormDataToRequest(data));

        if (isAPIError(response)) {
            toast.error(`Unable to check: ${response.message}`);
        } else {
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
        watch,
    };
}

export const mapCheckFormDataToRequest = (data: CheckFormData): any => {
    return {
        checks: [
            {
                action: data.action.id || '',
                resource_kind: data.resource.raw?.kind || '',
                resource_value: data.resource.raw?.value || '',
                principal: data.principal.id || '',        
            },
        ],
    };
}