

#include <arm_neon.h>

// hamming only works with length >= 16
void hamming(float *a, float *b, float *res, long *len)
{
    int size = *len;

    // use the vectorized version for the first n - (n % 4) elements
    int l = size - (size % 4);

    // create 4*4 registers to store the result
    uint32x4_t res_vec0 = vdupq_n_u32(0);
    uint32x4_t res_vec1 = vdupq_n_u32(0);
    uint32x4_t res_vec2 = vdupq_n_u32(0);
    uint32x4_t res_vec3 = vdupq_n_u32(0);

    int i = 0;

    uint32x4_t imr_1 = vdupq_n_u32(0);
    uint32x4_t imr_2 = vdupq_n_u32(0);
    uint32x4_t imr_3 = vdupq_n_u32(0);
    uint32x4_t imr_4 = vdupq_n_u32(0);

    // load 4*4 floats at a time
    while (i + 16 <= l)
    {
        float32x4x4_t a4 = vld1q_f32_x4(a + i);
        float32x4x4_t b4 = vld1q_f32_x4(b + i);

        res_vec0 -= vreinterpretq_s32_f32(vceqq_f32(a4.val[0], b4.val[0]));
        res_vec1 -= vreinterpretq_s32_f32(vceqq_f32(a4.val[1], b4.val[1]));
        res_vec2 -= vreinterpretq_s32_f32(vceqq_f32(a4.val[2], b4.val[2]));
        res_vec3 -= vreinterpretq_s32_f32(vceqq_f32(a4.val[3], b4.val[3]));

        i += 16;
    }

    while (i < l)
    {
        float32x4_t a_vec = vld1q_f32(a + i);
        float32x4_t b_vec = vld1q_f32(b + i);
        res_vec0 -= vreinterpretq_s32_f32(vceqq_f32(a_vec, b_vec));

        i += 4;
    }

    // convert to f32 implicitly
    int32_t sum = size;
    sum -= vaddvq_u32(res_vec0);
    sum -= vaddvq_u32(res_vec1);
    sum -= vaddvq_u32(res_vec2);
    sum -= vaddvq_u32(res_vec3);

    // add the remaining vectors
    for (int i = l; i < size; i++)
    {
        if (a[i] == b[i])
        {
            sum--;
        }
    }

    res[0] = sum;
}