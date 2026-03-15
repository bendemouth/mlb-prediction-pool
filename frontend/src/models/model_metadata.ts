export interface ModelMetadata {
    model_id: string;
    model_name: string;
    user_id: string;
    file_name: string;
    s3_key: string;
    status: string;
    created_at: string | Date;
    updated_at: string | Date;
}