import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { object, string } from 'yup';

import { Control, DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue } from 'react-hook-form';
import { createResource, updateResource } from 'service/model/resource';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';
import { Resource, ResourceAttribute } from 'service/model/model';

export type ResourceFormData = {
    id: string
    kind: string
    value: string
    attributes?: ResourceAttribute[]
}

export type OnSubmitHandler = (token: string) => FormEventHandler;

type ResourceForm = {
    control: Control<ResourceFormData>
    onSubmit: OnSubmitHandler
    register: UseFormRegister<ResourceFormData>
    getValues: UseFormGetValues<ResourceFormData>
    setValue: UseFormSetValue<ResourceFormData>
    defaultValues?: Readonly<DeepPartial<ResourceFormData>> | ResourceFormData
    errors: FieldErrors<ResourceFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<ResourceFormData>
}

const schema = object({
    id: string().required('You have to specify a name.'),
    kind: string().required('You have to specify a kind of resource.'),
    value: string().required('You have to specify a resource value (or at least wildcard).'),
}).required();

export default function useResourceForm(
    navigate: NavigateFunction,
    resource?: Resource,
): ResourceForm {
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
    } = useForm<ResourceFormData>({
        resolver: yupResolver(schema),
        defaultValues: mapResourceToFormData(resource),
    });

    const onSubmit = (token: string) => handleSubmit(async (data: ResourceFormData) => {
        if (resource === undefined) {
            const response = await createResource(token, mapResourceFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to create resource: ${response.message}`);
            } else {
                toast.success(`Resource ${response.id} has been created.`);
                navigate('/resources');
            }
        } else {
            const response = await updateResource(token, resource.id, mapResourceFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to update resource: ${response.message}`);
            } else {
                toast.success(`Resource ${response.id} has been updated.`);
                navigate('/resources');
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

export const mapResourceToFormData = (data?: Resource): ResourceFormData => {
    return ({
        id: data?.id || '',
        kind: data?.kind || '',
        value: data?.value || '',
        attributes: data?.attributes ? data.attributes : [],
    });
}

export const mapResourceFormDataToRequest = (data: ResourceFormData): any => {
    return {
        id: data.id,
        kind: data.kind,
        value: data.value,
        attributes: data.attributes,
    };
}