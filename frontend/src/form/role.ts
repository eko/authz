import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { array, mixed, object, string } from 'yup';

import { DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue } from 'react-hook-form';
import { createRole, updateRole } from 'service/model/role';
import { Role } from 'service/model/model';
import { ItemType } from 'component/MultipleAutocompleteInput';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';

export type RoleFormData = {
    id: string
    policies: ItemType[]
}

export type OnSubmitHandler = (token: string) => FormEventHandler;

type RoleForm = {
    onSubmit: OnSubmitHandler
    register: UseFormRegister<RoleFormData>
    getValues: UseFormGetValues<RoleFormData>
    setValue: UseFormSetValue<RoleFormData>
    defaultValues?: Readonly<DeepPartial<RoleFormData>> | RoleFormData
    errors: FieldErrors<RoleFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<RoleFormData>
}

const schema = object({
    id: string().required('You have to specify a name.'),
    policies: array()
        .of(mixed<ItemType>().defined())
        .min(1, 'You have to pick at least one policy.')
        .required('This field is required.'),
}).required();

export default function useRoleForm(
    navigate: NavigateFunction,
    role?: Role,
): RoleForm {
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
    } = useForm<RoleFormData>({
        resolver: yupResolver(schema),
        defaultValues: mapRoleToFormData(role),
    });

    const onSubmit = (token: string) => handleSubmit(async (data: RoleFormData) => {
        if (role === undefined) {
            const response = await createRole(token, mapRoleFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to create role: ${response.message}`);
            } else {
                toast.success(`Role ${response.id} has been created.`);
                navigate('/roles');
            }
        } else {
            const response = await updateRole(token, role.id, mapRoleFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to update role: ${response.message}`);
            } else {
                toast.success(`Role ${response.id} has been updated.`);
                navigate('/roles');
            }
        }
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

export const mapRoleToFormData = (role?: Role): RoleFormData => {
    return ({
        id: role?.id || '',
        policies: role?.policies.map(policy => ({
            id: policy.id,
            label: policy.id,
        }))  || [],
    })
}

export const mapRoleFormDataToRequest = (data: RoleFormData): any => {
    return {
        id: data.id,
        policies: data.policies.map(policy => policy.id) || [],
    };
}