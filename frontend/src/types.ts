export type BannerType = 'standard' | 'limited';

export interface BannerRequest {
  current_pity: number;
  planned_pulls: number;
  guaranteed: boolean;
}

export interface BannerResponse {
  total_5_star_probability: number;
  character_probability?: number;
  light_cone_probability?: number;
  rate_up_probability?: number;
  standard_char_probability?: number;
} 