import { FormEventHandler } from 'react';
import { yupResolver } from '@hookform/resolvers/yup';
import { array, mixed, object, string } from 'yup';

import { Control, DeepPartial, FieldErrors, useForm, UseFormGetValues, UseFormRegister, UseFormReset, UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { createPolicy, updatePolicy } from 'service/model/policy';
import { Policy } from 'service/model/model';
import { ItemType } from 'component/MultipleAutocompleteInput';
import { NavigateFunction } from 'react-router';
import { useToast } from 'context/toast';
import { isAPIError } from 'service/error/model';

export type PolicyFormData = {
    id: string
    resources: ItemType[]
    actions: ItemType[]
    attribute_rules?: string[]
}

export type OnSubmitHandler = (token: string) => FormEventHandler;

type PolicyForm = {
    control: Control<PolicyFormData>
    onSubmit: OnSubmitHandler
    register: UseFormRegister<PolicyFormData>
    getValues: UseFormGetValues<PolicyFormData>
    setValue: UseFormSetValue<PolicyFormData>
    defaultValues?: Readonly<DeepPartial<PolicyFormData>> | PolicyFormData
    errors: FieldErrors<PolicyFormData>
    isSubmitting: boolean
    isValid: boolean
    reset: UseFormReset<PolicyFormData>
    watch: UseFormWatch<PolicyFormData>
}

const schema = object({
    id: string().required('You have to specify a name.'),
    resources: array()
        .of(mixed<ItemType>())
        .min(1, 'You have to pick at least one policy.')
        .required('This field is required.'),
    actions: array()
        .of(mixed<ItemType>())
        .min(1, 'You have to pick at least one action.')
        .required('This field is required.'),
}).required();

export default function usePolicyForm(
    navigate: NavigateFunction,
    policy?: Policy,
): PolicyForm {
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
    } = useForm<PolicyFormData>({
        resolver: yupResolver(schema),
        defaultValues: mapPolicyToFormData(policy),
    });

    const onSubmit = (token: string) => handleSubmit(async (data: PolicyFormData) => {
        if (policy === undefined) {
            const response = await createPolicy(token, mapPolicyFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to create policy: ${response.message}`);
            } else {
                toast.success(`Policy ${response.id} has been created.`);
                navigate('/policies');
            }
        } else {
            const response = await updatePolicy(token, policy.id, mapPolicyFormDataToRequest(data));

            if (isAPIError(response)) {
                toast.error(`Unable to update policy: ${response.message}`);
            } else {
                toast.success(`Policy ${response.id} has been updated.`);
                navigate('/policies');
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
        watch,
    };
}

export const mapPolicyToFormData = (policy?: Policy): PolicyFormData => {
    return ({
        id: policy?.id || '',
        resources: policy?.resources.map(resource => ({
            id: resource.id,
            label: resource.id,
        }))  || [],
        actions: policy?.actions.map(action => ({
            id: action.id,
            label: action.id,
        }))  || [],
        attribute_rules: policy?.attribute_rules || [],
    })
}

export const mapPolicyFormDataToRequest = (data: PolicyFormData): any => {
    return {
        id: data.id,
        resources: data.resources.map(resource => resource.id) || [],
        actions: data.actions.map(action => action.id) || [],
        attribute_rules: data.attribute_rules || [],
    };
}