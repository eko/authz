import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { array, mixed, object, string } from 'yup';

import { Control, DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue } from 'react-hook-form';
import { createPrincipal, updatePrincipal } from 'service/model/principal';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';
import { Principal, PrincipalAttribute } from 'service/model/model';
import { ItemType } from 'component/MultipleAutocompleteInput';

export type PrincipalFormData = {
    id: string
    roles: ItemType[]
    attributes?: PrincipalAttribute[]
}

export type OnSubmitHandler = (token: string) => FormEventHandler;

type PrincipalForm = {
    control: Control<PrincipalFormData>
    onSubmit: OnSubmitHandler
    register: UseFormRegister<PrincipalFormData>
    getValues: UseFormGetValues<PrincipalFormData>
    setValue: UseFormSetValue<PrincipalFormData>
    defaultValues?: Readonly<DeepPartial<PrincipalFormData>> | PrincipalFormData
    errors: FieldErrors<PrincipalFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<PrincipalFormData>
}

const schema = object({
    id: string().required('You have to specify a name.'),
    roles: array().of(mixed<ItemType>()),
}).required();

export default function usePrincipalForm(
    navigate: NavigateFunction,
    principal?: Principal,
): PrincipalForm {
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
    } = useForm<PrincipalFormData>({
        resolver: yupResolver(schema),
        defaultValues: mapPrincipalToFormData(principal),
    });

    const onSubmit = (token: string) => handleSubmit(async (data: PrincipalFormData) => {
        if (principal === undefined) {
            const response = await createPrincipal(token, mapPrincipalFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to create principal: ${response.message}`);
            } else {
                toast.success(`Principal ${response.id} has been created.`);
                navigate('/principals');
            }
        } else {
            const response = await updatePrincipal(token, principal.id, mapPrincipalFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to update principal: ${response.message}`);
            } else {
                toast.success(`Principal ${response.id} has been updated.`);
                navigate('/principals');
            }
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

export const mapPrincipalToFormData = (data?: Principal): PrincipalFormData => {
    return ({
        id: data?.id || '',
        roles: data?.roles?.map(role => ({
            id: role.id,
            label: role.id,
        }))  || [],
        attributes: data?.attributes ? data.attributes : [],
    });
}

export const mapPrincipalFormDataToRequest = (data: PrincipalFormData): any => {
    return {
        id: data.id,
        roles: data?.roles.map(role => role.id)  || [],
        attributes: data.attributes,
    };
}